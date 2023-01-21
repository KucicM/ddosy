package ddosy_test

import (
	"testing"

	ddosy "github.com/kucicm/ddosy/app"
)

func TestScheduleTask(t *testing.T) {
	req := ddosy.ScheduleRequestWeb{}

	repo := ddosy.NewTaskRepository("test.db", true)
	provider := ddosy.NewTaskProvider(repo)
	id, err := provider.ScheduleTask(req)
	if err != nil {
		t.Errorf("got error %s", err)
	}

	if id != 1 {
		t.Errorf("expected id == 1 got %d", id)
	}
}

func TestGetTaskAfterScheduled(t *testing.T) {
	req := ddosy.ScheduleRequestWeb{}

	repo := ddosy.NewTaskRepository("test.db", true)
	provider := ddosy.NewTaskProvider(repo)
	provider.ScheduleTask(req)

	task := provider.Next()
	if task == nil {
		t.Error("got no task")
	}
}
