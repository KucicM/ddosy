package ddosy_test

import (
	"testing"
	"time"

	ddosy "github.com/kucicm/ddosy/app"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func TestSinglePatternMetrics(t *testing.T) {
	var id uint64 = 10
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

	provider := ddosy.NewRelustProvider()
	provider.NewPattern(id)

	provider.Update(id, res)

	provider.Done(id)

	out, err := provider.Get(id)
	if err != nil {
		t.Error(err)
	}

	if len(out) < 400 {
		t.Errorf("unexpected result, got %s\n", out)
	}
}

func TestTwoPatternsMetrics(t *testing.T) {
	var id uint64 = 10
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

	provider := ddosy.NewRelustProvider()

	provider.NewPattern(id)
	provider.Update(id, res)

	provider.NewPattern(id)
	provider.Update(id, res)

	provider.Done(id)

	out, err := provider.Get(id)
	if err != nil {
		t.Error(err)
	}

	if len(out) < 800 {
		t.Errorf("unexpected result, got %s\n", out)
	}
}
