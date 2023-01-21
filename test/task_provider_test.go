package ddosy_test

import (
	"testing"

	ddosy "github.com/kucicm/ddosy/app"
)

func TestScheduleTask(t *testing.T) {
	req := ddosy.ScheduleRequestWeb{}

	repo := ddosy.NewTaskRepository(":memory:")
	provider := ddosy.NewTaskProvider(repo)
	id, err := provider.ScheduleTask(req)
	if err != nil {
		t.Errorf("got error %s", err)
	}

	if id != 1 {
		t.Errorf("expected id == 1 got %d", id)
	}
}
