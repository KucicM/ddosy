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

func TestStatusUpdateFlow(t *testing.T) {
	rep := ddosy.NewTaskRepository(":memory:")

	req := ddosy.ScheduleRequestWeb{
		Endpoint:        "test-url",
		TrafficPatterns: []ddosy.TrafficPatternWeb{{Weight: 1, Payload: "test"}},
		LoadPatterns: []ddosy.LoadPatternWeb{{
			Duration: "1s",
			Linear:   &ddosy.LinearLoadWeb{StartRate: 50, EndRate: 50},
		}},
	}

	id, _ := rep.Save(req)
	scheduled, _ := rep.Get(id)

	// scheduled to running
	if err := rep.UpdateStatus(id, ddosy.Running); err != nil {
		t.Errorf("error in update %s\n", err)
	}

	running, err := rep.Get(id)
	if err != nil {
		t.Errorf("unexpected error on get %s\n", err)
	}

	if running.StatusId != ddosy.Running {
		t.Errorf("expected %d got %d\n", ddosy.Running, running.StatusId)
	}

	if running.CreatedAt.UnixMilli() != scheduled.CreatedAt.UnixMilli() {
		t.Errorf("unexpected update on createdAt column expected %v got %v\n", scheduled.CreatedAt, running.CreatedAt)
	}

	if running.StartedAt == nil {
		t.Error("expecting update at startedAt")
	}

	if running.KilledAt != nil {
		t.Error("unexpected update on KilledAt column")
	}

	if running.DoneAt != nil {
		t.Error("unexpected update on DoneAt column")
	}

	// running to done
	if err := rep.UpdateStatus(id, ddosy.Done); err != nil {
		t.Errorf("error in update %s\n", err)
	}

	done, err := rep.Get(id)
	if err != nil {
		t.Errorf("unexpected error on get %s\n", err)
	}

	if done.StatusId != ddosy.Done {
		t.Errorf("expected %d got %d\n", ddosy.Done, done.StatusId)
	}

	if done.CreatedAt.UnixMilli() != running.CreatedAt.UnixMilli() {
		t.Errorf("unexpected update on createdAt column expected %v got %v\n", running.CreatedAt, done.CreatedAt)
	}

	if done.StartedAt.UnixMilli() != running.StartedAt.UnixMilli() {
		t.Errorf("unexpected update on StartedAt column expected %v got %v\n", running.StartedAt, done.StartedAt)
	}

	if done.KilledAt != nil {
		t.Error("unexpected update on KilledAt column")
	}

	if done.DoneAt == nil {
		t.Error("expecting update at DoneAt")
	}

	// done -> killed?
	if err := rep.UpdateStatus(id, ddosy.Killed); err != nil {
		t.Errorf("error in update %s\n", err)
	}

	killed, err := rep.Get(id)
	if err != nil {
		t.Errorf("unexpected error on get %s\n", err)
	}

	if !reflect.DeepEqual(killed, done) {
		t.Errorf("expected no change %+v != %+v\n", killed, done)
	}
}

func TestRunningToKilledStatusUpdate(t *testing.T) {
	rep := ddosy.NewTaskRepository(":memory:")
	id, _ := rep.Save(ddosy.ScheduleRequestWeb{})
	rep.UpdateStatus(id, ddosy.Running)
	running, _ := rep.Get(id)

	err := rep.UpdateStatus(id, ddosy.Killed)
	if err != nil {
		t.Errorf("unexpected error %s\n", err)
	}

	killed, err := rep.Get(id)
	if err != nil {
		t.Errorf("unexpected error on get %s\n", err)
	}

	if killed.StatusId != ddosy.Killed {
		t.Errorf("expected %d got %d\n", ddosy.Killed, killed.StatusId)
	}

	if killed.CreatedAt.UnixMilli() != running.CreatedAt.UnixMilli() {
		t.Errorf("unexpected update on createdAt column expected %v got %v\n", running.CreatedAt, killed.CreatedAt)
	}

	if killed.StartedAt.UnixMilli() != running.StartedAt.UnixMilli() {
		t.Errorf("unexpected update on StartedAt column expected %v got %v\n", running.StartedAt, killed.StartedAt)
	}

	if killed.KilledAt == nil {
		t.Error("expecting update at KilledAt")
	}

	if killed.DoneAt != nil {
		t.Error("unexpected update on DoneAt column")
	}

	// killed -> done?
	if err := rep.UpdateStatus(id, ddosy.Done); err != nil {
		t.Errorf("error in update %s\n", err)
	}

	done, err := rep.Get(id)
	if err != nil {
		t.Errorf("unexpected error on get %s\n", err)
	}

	if !reflect.DeepEqual(killed, done) {
		t.Errorf("expected no change %+v != %+v\n", killed, done)
	}
}