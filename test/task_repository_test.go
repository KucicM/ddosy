package ddosy_test

import (
	"reflect"
	"testing"
	"time"

	ddosy "github.com/kucicm/ddosy/app"
)

func TestSave(t *testing.T) {
	rep := ddosy.NewTaskRepository("test.db", true)

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
	rep := ddosy.NewTaskRepository("test.db", true)

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
	rep := ddosy.NewTaskRepository("test.db", true)

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
	rep := ddosy.NewTaskRepository("test.db", true)
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

func TestUpdateProgress(t *testing.T) {
	rep := ddosy.NewTaskRepository("test.db", true)
	id, _ := rep.Save(ddosy.ScheduleRequestWeb{})

	// 1
	if err := rep.UpdateProgress(id, "test"); err != nil {
		t.Errorf("unexpected error %s\n", err)
	}

	o1, err := rep.Get(id)
	if err != nil {
		t.Errorf("unexpected error %s\n", err)
	}

	if o1.Results != "\ntest" {
		t.Errorf("unexpected relusts %s\n", o1.Results)
	}

	// 2
	if err := rep.UpdateProgress(id, "test2"); err != nil {
		t.Errorf("unexpected error %s\n", err)
	}

	o1, err = rep.Get(id)
	if err != nil {
		t.Errorf("unexpected error %s\n", err)
	}

	if o1.Results != "\ntest\ntest2" {
		t.Errorf("unexpected relusts %s\n", o1.Results)
	}
}

func TestNoTaskGetNext(t *testing.T) {
	rep := ddosy.NewTaskRepository("test.db", true)
	id, _, err := rep.GetNextTask()
	if err != nil {
		t.Errorf("unexpected error %s\n", err)
	}

	if id != 0 {
		t.Errorf("got id = %d expcted 0\n", id)
	}
}

func TestGetNextTask(t *testing.T) {
	rep := ddosy.NewTaskRepository("test.db", true)

	killedId, _ := rep.Save(ddosy.ScheduleRequestWeb{Endpoint: "1"})
	rep.UpdateStatus(killedId, ddosy.Killed)

	doneId, _ := rep.Save(ddosy.ScheduleRequestWeb{Endpoint: "2"})
	rep.UpdateStatus(doneId, ddosy.Running)
	rep.UpdateStatus(doneId, ddosy.Done)

	// check if no value is returned
	id, _, err := rep.GetNextTask()
	if err != nil {
		t.Errorf("unexpected error %s\n", err)
	}

	if id != 0 {
		t.Errorf("got id = %d expcted 0\n", id)
	}

	scheduled1, _ := rep.Save(ddosy.ScheduleRequestWeb{Endpoint: "3"})
	scheduled2, _ := rep.Save(ddosy.ScheduleRequestWeb{Endpoint: "4"})

	if scheduled1 == scheduled2 {
		t.Error("why are two ids equal?")
	}

	id, req, err := rep.GetNextTask()
	if err != nil {
		t.Errorf("unexpected error %s\n", err)
	}

	if id != scheduled1 {
		t.Errorf("expected id=%d got %d\n", scheduled1, id)
	}

	if req.Endpoint != "3" {
		t.Errorf("did not get right request %v", req)
	}

}
