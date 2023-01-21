package ddosy_test

import (
	"reflect"
	"testing"
	"time"

	ddosy "github.com/kucicm/ddosy/app"
)

func TestRepositoryBasic(t *testing.T) {
	// rep := ddosy.NewTaskRepository(":memory:")
	// log.Println(rep.InsertNew(ddosy.LoadTask{}))
	// log.Println(rep.InsertNew(ddosy.LoadTask{}))
	// log.Println(rep.InsertNew(ddosy.LoadTask{}))
	// rep.Close()
}

func TestSave(t *testing.T) {
	rep := ddosy.NewTaskRepository(":memory:")

	req := ddosy.ScheduleRequestWeb{
		Endpoint:        "test-url",
		TrafficPatterns: []ddosy.TrafficPatternWeb{{Weight: 1, Payload: "test"}},
		LoadPatterns: []ddosy.LoadPatternWeb{{
			Duration: "1s",
			Linear:   &ddosy.LinearLoadWeb{StartRate: 50, EndRate: 50},
		}},
	}

	id, err := rep.Save(req)
	if err != nil {
		t.Errorf("unexpected error on save %s\n", err)
	}

	if id != 1 {
		t.Errorf("expected id=1 got id=%d\n", id)
	}

}

func TestGet(t *testing.T) {
	rep := ddosy.NewTaskRepository(":memory:")

	req := ddosy.ScheduleRequestWeb{
		Endpoint:        "test-url",
		TrafficPatterns: []ddosy.TrafficPatternWeb{{Weight: 1, Payload: "test"}},
		LoadPatterns: []ddosy.LoadPatternWeb{{
			Duration: "1s",
			Linear:   &ddosy.LinearLoadWeb{StartRate: 50, EndRate: 50},
		}},
	}

	id, err := rep.Save(req)
	if err != nil {
		t.Errorf("unexpected error on save %s\n", err)
	}

	task, err := rep.Get(id)
	if err != nil {
		t.Errorf("unexpected error on get %s\n", err)
	}

	if task.Id != id {
		t.Errorf("expected id=%d got id=%d\n", id, task.Id)
	}

	if task.CreatedAt.After(time.Now()) {
		t.Errorf("time is going backwards? task.CreatedAt=%s\n", task.CreatedAt)
	}

	if task.StartedAt != nil {
		t.Errorf("expected nil got %s\n", task.StartedAt)
	}

	if task.DoneAt != nil {
		t.Errorf("expected nil got %s\n", task.DoneAt)
	}

	if !reflect.DeepEqual(task.Request, req) {
		t.Errorf("expected %+v got %+v\n", req, task.Request)
	}

	if task.Results != "" {
		t.Errorf("expected empty string got %s\n", task.Results)
	}
}

func TestStatusUpdate(t *testing.T) {
	
}