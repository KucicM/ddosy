package ddosy_test

// import (
// 	"testing"
// 	ddosy "github.com/kucicm/ddosy/pkg"
// )

// func TestValidTrafficPattern(t *testing.T) {
// 	pattern := ddosy.LoadPattern
// }

// func TestValidConstantLoadPattern(t *testing.T) {
// 	str := `{
// 		"duration": "1m",
// 		"linear": {
// 			"startRate": 10,
// 			"endRate": 10
// 		}
// 	}`

// 	var pattern ddosy.LoadPattern
// 	if err := json.Unmarshal([]byte(str), &pattern); err != nil {
// 		t.Errorf("unable to unmarshal load pattern %s\n", err)
// 	}

// 	if err := pattern.Validate(); err != nil {
// 		t.Errorf("validation failed %s\n", err)
// 	}

// 	pacer, err := pattern.Pacer()
// 	if err != nil {
// 		t.Errorf("failed to get pacer %s\n", err)
// 	}

// 	expected := vegeta.LinearPacer{
// 		StartAt: vegeta.Rate{Freq: 10, Per: time.Second},
// 		Slope:   0,
// 	}
// 	if !reflect.DeepEqual(pacer, expected) {
// 		t.Errorf("Expected %+v got %+v\n", expected, pacer)
// 	}
// }

// func TestValidLinearLoadPattern(t *testing.T) {
// 	str := `{
// 		"duration": "1s",
// 		"linear": {
// 			"startRate": 10,
// 			"endRate": 100
// 		}
// 	}`

// 	var pattern ddosy.LoadPattern
// 	if err := json.Unmarshal([]byte(str), &pattern); err != nil {
// 		t.Errorf("unable to unmarshal load pattern %s\n", err)
// 	}

// 	if err := pattern.Validate(); err != nil {
// 		t.Errorf("validation failed %s\n", err)
// 	}

// 	pacer, err := pattern.Pacer()
// 	if err != nil {
// 		t.Errorf("failed to get pacer %s\n", err)
// 	}

// 	expected := vegeta.LinearPacer{
// 		StartAt: vegeta.Rate{Freq: 10, Per: time.Second},
// 		Slope:   90,
// 	}
// 	if !reflect.DeepEqual(pacer, expected) {
// 		t.Errorf("Expected %+v got %+v\n", expected, pacer)
// 	}
// }

// func TestInvalidDurationLinearLoadPattern(t *testing.T) {
// 	str := `{
// 		"duration": "1",
// 		"linear": {
// 			"startRate": 10,
// 			"endRate": 100
// 		}
// 	}`

// 	var pattern ddosy.LoadPattern
// 	if err := json.Unmarshal([]byte(str), &pattern); err != nil {
// 		t.Errorf("unable to unmarshal load pattern %s\n", err)
// 	}

// 	if err := pattern.Validate(); err == nil {
// 		t.Error("duration is invalid but validation did not fail")
// 	}
// }

// func TestNegativeDurationLinearLoadPattern(t *testing.T) {
// 	str := `{
// 		"duration": "-1m",
// 		"linear": {
// 			"startRate": 10,
// 			"endRate": 100
// 		}
// 	}`

// 	var pattern ddosy.LoadPattern
// 	if err := json.Unmarshal([]byte(str), &pattern); err != nil {
// 		t.Errorf("unable to unmarshal load pattern %s\n", err)
// 	}

// 	if err := pattern.Validate(); err == nil {
// 		t.Error("duration is negative but validation did not fail")
// 	}
// }
