package ddosy_test

import (
	"testing"

	ddosy "github.com/kucicm/ddosy/pkg"
)


func TestScheduleTask(t *testing.T) {
	task := &ddosy.LoadTask{}

	provider := ddosy.NewTaskProvider(5)
	id, err := provider.ScheduleTask(task)
	if err != nil {
		t.Errorf("got error %s", err)
	}

	if id != 1 {
		t.Errorf("expected id == 1 got %d", id)
	}
}

func TestMaxQueue(t *testing.T) {
	task := &ddosy.LoadTask{}
	provider := ddosy.NewTaskProvider(2)

	id, err := provider.ScheduleTask(task)
	if err != nil {
		t.Errorf("got error %s", err)
	}

	if id != 1 {
		t.Errorf("expected id == 1 got %d", id)
	}

	id, err = provider.ScheduleTask(task)
	if err != nil {
		t.Errorf("got error %s", err)
	}

	if id != 2 {
		t.Errorf("expected id == 1 got %d", id)
	}

	if _, err = provider.ScheduleTask(task); err == nil {
		t.Error("expected error got nil")
	}
}