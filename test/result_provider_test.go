package ddosy_test

import (
	"log"
	"testing"
	"time"

	ddosy "github.com/kucicm/ddosy/app"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func TestSinglePatternMetrics(t *testing.T) {
	res := &vegeta.Result{
		Attack:    "attack",
		Seq:       1,
		Code:      200,
		Timestamp: time.Now(),
		Latency:   time.Millisecond,
		BytesOut:  30,
		BytesIn:   40,
		Error:     "",
		Body:      nil,
		Method:    "GET",
		URL:       "URL",
		Headers:   nil,
	}

	repo := ddosy.NewTaskRepository("test.db", true)
	id, _ := repo.Save(ddosy.ScheduleRequestWeb{Endpoint: "test"})
	provider := ddosy.NewRelustProvider(repo)
	provider.UpdateRunning(id, res)
	provider.FinalizeRunning(id)

	out, err := provider.Get(id)
	log.Println(out)
	if err != nil {
		t.Error(err)
	}

	if len(out) < 400 {
		t.Errorf("unexpected result, got %s\n", out)
	}
}

func TestTwoPatternsMetrics(t *testing.T) {
	res := &vegeta.Result{
		Attack:    "attack",
		Seq:       1,
		Code:      200,
		Timestamp: time.Now(),
		Latency:   time.Millisecond,
		BytesOut:  30,
		BytesIn:   40,
		Error:     "",
		Body:      nil,
		Method:    "GET",
		URL:       "URL",
		Headers:   nil,
	}

	repo := ddosy.NewTaskRepository("test.db", true)
	provider := ddosy.NewRelustProvider(repo)

	id, _ := repo.Save(ddosy.ScheduleRequestWeb{Endpoint: "test"})

	provider.UpdateRunning(id, res)
	provider.FinalizeRunning(id)
	provider.UpdateRunning(id, res)

	out, err := provider.Get(id)
	log.Println(out)
	if err != nil {
		t.Error(err)
	}

	if len(out) < 800 {
		t.Errorf("unexpected result, got %s\n", out)
	}
}
