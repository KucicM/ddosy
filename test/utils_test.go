package ddosy_test

import (
	"reflect"
	"testing"

	ddosy "github.com/kucicm/ddosy/app"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func TestNormDistribution(t *testing.T) {
	dist := []float64{1, 2, 1}
	normDist := ddosy.CreateDistribution(dist)

	expected := []float64{0.25, 0.75, 1}
	for i, actual := range normDist {
		if actual != expected[i] {
			t.Errorf("expected %+v got %+v", expected, actual)
		}
	}
}

func TestSinglePattern(t *testing.T) {
	p := ddosy.TrafficPatternWeb{
		Weight:  10,
		Payload: "haha",
	}
	endpoint := "test-endpoint"
	pattern := ddosy.NewTrafficPattern(endpoint, []ddosy.TrafficPatternWeb{p})

	targeter := ddosy.NewWeightedTargeter(pattern)

	target := &vegeta.Target{}
	if err := targeter(target); err != nil {
		t.Errorf("got error from targeter %s", err)
	}

	expectedBody := []byte(p.Payload)
	if !reflect.DeepEqual(expectedBody, target.Body) {
		t.Errorf("expected %v got %v", expectedBody, target.Body)
	}

	if endpoint != target.URL {
		t.Errorf("expected %v got %v", endpoint, target.URL)
	}
}

func TestMultPattern(t *testing.T) {
	p1 := ddosy.TrafficPatternWeb{
		Weight:  10,
		Payload: "1",
	}
	p2 := ddosy.TrafficPatternWeb{
		Weight:  10,
		Payload: "2",
	}
	endpoint := "test-endpoint"
	pattern := ddosy.NewTrafficPattern(endpoint, []ddosy.TrafficPatternWeb{p1, p2})

	targeter := ddosy.NewWeightedTargeter(pattern)

	numOfIters := 1000
	tracker := make(map[string]int)
	for i := 0; i < numOfIters; i++ {
		target := &vegeta.Target{}
		targeter(target)
		tracker[string(target.Body)] += 1
	}

	if len(tracker) != 2 {
		t.Errorf("something dose not work with targeter, tracker size %+v", tracker)
	}

	for _, v := range tracker {
		if v < 100 {
			t.Errorf("Expected 50%% distribution got %d%%", v*100/numOfIters)
		}
	}
}
