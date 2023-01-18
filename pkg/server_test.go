package ddosy_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	ddosy "github.com/kucicm/ddosy/pkg"
)


func TestRunLoadTest(t *testing.T) {
	cfg := ddosy.ServerConfig{Port: 434343, MaxQueue: 5}
	srv := ddosy.NewServer(cfg)

	sr := ddosy.ScheduleRequestWeb{
		Endpoint: "localdost:777777/test",
		TrafficPatterns: []ddosy.TrafficPatternWeb{{Weight: 1, Payload: "test"}},
		LoadPatterns: []ddosy.LoadPatternWeb{{
			Duration: "1s", 
			Linear: &ddosy.LinearLoadWeb{StartRate: 1, EndRate: 1},
		}},
	}

    var buf bytes.Buffer
    json.NewEncoder(&buf).Encode(sr)

	req := httptest.NewRequest(http.MethodPost, "/run", &buf)
    w := httptest.NewRecorder()

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
}