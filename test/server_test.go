package ddosy_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	ddosy "github.com/kucicm/ddosy/app"
)

func TestRunLoadTest(t *testing.T) {
	var callCount int32
	done := make(chan struct{})
	testSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&callCount, 1) >= 5 {
			done <- struct{}{}
		}
	}))
	defer testSrv.Close()

	sr := ddosy.ScheduleRequestWeb{
		Endpoint:        testSrv.URL,
		TrafficPatterns: []ddosy.TrafficPatternWeb{{Weight: 1, Payload: "test"}},
		LoadPatterns: []ddosy.LoadPatternWeb{{
			Duration: "1s",
			Linear:   &ddosy.LinearLoadWeb{StartRate: 50, EndRate: 50},
		}},
	}

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(sr)
	req := httptest.NewRequest(http.MethodPost, "/run", &buf)
	w := httptest.NewRecorder()

	cfg := ddosy.ServerConfig{Port: 434343, DbUrl: "test.db", TruncateDbOnStart: true}
	srv := ddosy.NewServer(cfg)
	srv.ScheduleHandler(w, req)

	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	var actual ddosy.ScheduleResponseWeb
	if err := json.Unmarshal(data, &actual); err != nil {
		t.Errorf("falied to unmarshal schedule response %s", err)
	}

	expected := ddosy.ScheduleResponseWeb{Id: 1}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected %v got %v", expected, actual)
	}

	select {
	case <-done:
	case <-time.After(time.Second * 3):
		t.Errorf("load test did not start %d", atomic.LoadInt32(&callCount))
	}
}

func TestKillRunningLoadTest(t *testing.T) {
	var callCount int32
	started := make(chan struct{})
	done := make(chan struct{})
	testSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		i := atomic.AddInt32(&callCount, 1)
		if i == 5 {
			started <- struct{}{}
		}
		if i == 50 {
			// should never reach
			done <- struct{}{}
		}
	}))
	defer testSrv.Close()

	sr := ddosy.ScheduleRequestWeb{
		Endpoint:        testSrv.URL,
		TrafficPatterns: []ddosy.TrafficPatternWeb{{Weight: 1, Payload: "test"}},
		LoadPatterns: []ddosy.LoadPatternWeb{{
			Duration: "10s",
			Linear:   &ddosy.LinearLoadWeb{StartRate: 50, EndRate: 50},
		}},
	}

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(sr)
	req := httptest.NewRequest(http.MethodPost, "/run", &buf)
	w := httptest.NewRecorder()

	cfg := ddosy.ServerConfig{Port: 434343, DbUrl: "test.db", TruncateDbOnStart: true}
	srv := ddosy.NewServer(cfg)
	srv.ScheduleHandler(w, req)

	select {
	case <-started:
	case <-time.After(time.Second * 3):
		t.Errorf("load test did not start %d\n", callCount)
	}

	req = httptest.NewRequest(http.MethodDelete, "/kill", nil)
	w = httptest.NewRecorder()
	srv.KillHandler(w, req)

	select {
	case <-done:
		t.Error("kill did not work")
	case <-time.After(time.Second * 3):
	}
}

func TestGetResults(t *testing.T) {
	var callCount int32
	started := make(chan struct{})
	testSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		i := atomic.AddInt32(&callCount, 1)
		if i == 5 {
			started <- struct{}{}
		}
	}))
	defer testSrv.Close()

	sr := ddosy.ScheduleRequestWeb{
		Endpoint:        testSrv.URL,
		TrafficPatterns: []ddosy.TrafficPatternWeb{{Weight: 1, Payload: "test"}},
		LoadPatterns: []ddosy.LoadPatternWeb{{
			Duration: "10s",
			Linear:   &ddosy.LinearLoadWeb{StartRate: 50, EndRate: 50},
		}},
	}

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(sr)
	req := httptest.NewRequest(http.MethodPost, "/run", &buf)
	w := httptest.NewRecorder()

	cfg := ddosy.ServerConfig{Port: 434343, DbUrl: "test.db", TruncateDbOnStart: true}
	srv := ddosy.NewServer(cfg)
	srv.ScheduleHandler(w, req)

	res := w.Result()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	res.Body.Close()

	var resp ddosy.ScheduleResponseWeb
	if err := json.Unmarshal(data, &resp); err != nil {
		t.Errorf("falied to unmarshal schedule response %s", err)
	}

	if resp.Error != "" {
		t.Errorf("got error from server %s\n", resp.Error)
	}

	select {
	case <-started:
	case <-time.After(time.Second * 3):
		t.Errorf("load test did not start %d\n", callCount)
	}

	req = httptest.NewRequest(http.MethodGet, "/satus", nil)
	q := req.URL.Query()
	q.Add("id", strconv.FormatUint(resp.Id, 10))
	req.URL.RawQuery = q.Encode()
	w = httptest.NewRecorder()
	srv.StatusHandler(w, req)
	res = w.Result()
	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	res.Body.Close()
	if len(data) < 100 {
		t.Errorf("unexpected result %s\n", string(data))
	}

}
