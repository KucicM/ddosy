package ddosy_test

import (
	"reflect"
	"testing"

	ddosy "github.com/kucicm/ddosy/pkg"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func TestSinglePattern(t *testing.T) {
	req := ddosy.TrafficPatternWeb{
		Weight: 10,
		Payload: "haha",
	}
	endpoint := "test-endpoint"
	pattern := ddosy.NewTrafficPattern(endpoint, []ddosy.TrafficPatternWeb{req})

	targeter := ddosy.NewWeightedTargeter(pattern)

	target := &vegeta.Target{}
	if err := targeter(target); err != nil {
		t.Errorf("got error from targeter %s", err)
	}

	expectedBody := []byte(req.Payload)
	if !reflect.DeepEqual(expectedBody, target.Body) {
		t.Errorf("expected %v got %v", expectedBody, target.Body)
	}

	if endpoint != target.URL {
		t.Errorf("expected %v got %v", endpoint, target.URL)
	}

}
