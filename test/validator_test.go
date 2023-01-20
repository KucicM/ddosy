package ddosy_test

import (
	"testing"

	ddosy "github.com/kucicm/ddosy/app"
)

func TestValidTrafficPatternWeb(t *testing.T) {
	pattern := ddosy.TrafficPatternWeb{Weight: 1}

	if err := ddosy.ValidateTrafficPatternWeb(pattern); err != nil {
		t.Errorf("expected no error got %s\n", err)
	}
}

func TestZeroWeigthTrafficPatternWeb(t *testing.T) {
	pattern := ddosy.TrafficPatternWeb{Weight: 0}

	if err := ddosy.ValidateTrafficPatternWeb(pattern); err == nil {
		t.Error("expected error got nil")
	}
}

func TestNegativeWeigthTrafficPatternWeb(t *testing.T) {
	pattern := ddosy.TrafficPatternWeb{Weight: -1}

	if err := ddosy.ValidateTrafficPatternWeb(pattern); err == nil {
		t.Error("expected error got nil")
	}
}

func TestValidLinearLoadWeb(t *testing.T) {
	load := ddosy.LinearLoadWeb{
		StartRate: 0,
		EndRate: 10,
	}

	if err := ddosy.ValidateLinearLoadWeb(load); err != nil {
		t.Errorf("expected no error got %s\n", err)
	}
}

func TestNegativeStartRateLinearLoadWeb(t *testing.T) {
	load := ddosy.LinearLoadWeb{
		StartRate: -1,
		EndRate: 10,
	}

	if err := ddosy.ValidateLinearLoadWeb(load); err == nil {
		t.Error("expected error got nil")
	}
}

func TestNegativeEndRateLinearLoadWeb(t *testing.T) {
	load := ddosy.LinearLoadWeb{
		StartRate: 0,
		EndRate: -10,
	}

	if err := ddosy.ValidateLinearLoadWeb(load); err == nil {
		t.Error("expected error got nil")
	}
}

func TestValidSineLoadWeb(t *testing.T) {
	load := ddosy.SineLoadWeb{
		Mean: 10,
		Amp: 5,
		Period: "1m",
	}

	if err := ddosy.ValidateSineLoadWeb(load); err != nil {
		t.Errorf("expected no error got %s\n", err)
	}
}

func TestEmptyPeriodSineLoadWeb(t *testing.T) {
	load := ddosy.SineLoadWeb{
		Mean: 10,
		Amp: 5,
		Period: "",
	}

	if err := ddosy.ValidateSineLoadWeb(load); err == nil {
		t.Error("expected error got nil")
	}
}

func TestNegativePeriodSineLoadWeb(t *testing.T) {
	load := ddosy.SineLoadWeb{
		Mean: 10,
		Amp: 5,
		Period: "-1m",
	}

	if err := ddosy.ValidateSineLoadWeb(load); err == nil {
		t.Error("expected error got nil")
	}
}

func TestZeroMeanSineLoadWeb(t *testing.T) {
	load := ddosy.SineLoadWeb{
		Mean: 0,
		Amp: 5,
		Period: "1m",
	}

	if err := ddosy.ValidateSineLoadWeb(load); err == nil {
		t.Error("expected error got nil")
	}
}

func TestNegativeApmSineLoadWeb(t *testing.T) {
	load := ddosy.SineLoadWeb{
		Mean: 10,
		Amp: -1,
		Period: "1m",
	}

	if err := ddosy.ValidateSineLoadWeb(load); err == nil {
		t.Error("expected error got nil")
	}
}

func TestLargeApmSineLoadWeb(t *testing.T) {
	load := ddosy.SineLoadWeb{
		Mean: 10,
		Amp: 10,
		Period: "1m",
	}

	if err := ddosy.ValidateSineLoadWeb(load); err == nil {
		t.Error("expected error got nil")
	}
}


func TestValidLoadPatternWeb(t *testing.T) {
	pattern := ddosy.LoadPatternWeb{
		Duration: "10m",
		Linear: &ddosy.LinearLoadWeb{},
	}

	if err := ddosy.ValidateLoadPatternWeb(pattern); err != nil {
		t.Errorf("expected no error got %s\n", err)
	}
}

func TestEmptyDurationLoadPatternWeb(t *testing.T) {
	pattern := ddosy.LoadPatternWeb{
		Duration: "",
		Linear: &ddosy.LinearLoadWeb{},
	}

	if err := ddosy.ValidateLoadPatternWeb(pattern); err == nil {
		t.Error("expected error got nil")
	}
}

func TestNegativeDurationLoadPatternWeb(t *testing.T) {
	pattern := ddosy.LoadPatternWeb{
		Duration: "-1m",
		Linear: &ddosy.LinearLoadWeb{},
	}

	if err := ddosy.ValidateLoadPatternWeb(pattern); err == nil {
		t.Error("expected error got nil")
	}
}

func TestNoPatternDefineLoadPatternWeb(t *testing.T) {
	pattern := ddosy.LoadPatternWeb{
		Duration: "1m",
	}

	if err := ddosy.ValidateLoadPatternWeb(pattern); err == nil {
		t.Error("expected error got nil")
	}
}

func TestInvalidLinearPatternDefineLoadPatternWeb(t *testing.T) {
	pattern := ddosy.LoadPatternWeb{
		Duration: "1m",
		Linear: &ddosy.LinearLoadWeb{StartRate: -1},
	}

	if err := ddosy.ValidateLoadPatternWeb(pattern); err == nil {
		t.Error("expected error got nil")
	}
}

func TestInvalidSinePatternDefineLoadPatternWeb(t *testing.T) {
	pattern := ddosy.LoadPatternWeb{
		Duration: "1m",
		Sine: &ddosy.SineLoadWeb{},
	}

	if err := ddosy.ValidateLoadPatternWeb(pattern); err == nil {
		t.Error("expected error got nil")
	}
}

func TestValidScheduleRequestWeb(t *testing.T) {
	load := ddosy.LoadPatternWeb{
		Duration: "10m",
		Linear: &ddosy.LinearLoadWeb{},
	}
	traffic := ddosy.TrafficPatternWeb{Weight: 1}

	req := ddosy.ScheduleRequestWeb{
		Endpoint: "test",
		LoadPatterns: []ddosy.LoadPatternWeb{load},
		TrafficPatterns: []ddosy.TrafficPatternWeb{traffic},
	}

	if err := ddosy.ValidateScheduleRequestWeb(req); err != nil {
		t.Errorf("expected no error got %s\n", err)
	}
}

func TestNoEndpointScheduleRequestWeb(t *testing.T) {
	load := ddosy.LoadPatternWeb{
		Duration: "10m",
		Linear: &ddosy.LinearLoadWeb{},
	}
	traffic := ddosy.TrafficPatternWeb{Weight: 1}

	req := ddosy.ScheduleRequestWeb{
		LoadPatterns: []ddosy.LoadPatternWeb{load},
		TrafficPatterns: []ddosy.TrafficPatternWeb{traffic},
	}

	if err := ddosy.ValidateScheduleRequestWeb(req); err == nil {
		t.Error("expected error got nil")
	}
}

func TestNoTrafficPatternsScheduleRequestWeb(t *testing.T) {
	load := ddosy.LoadPatternWeb{
		Duration: "10m",
		Linear: &ddosy.LinearLoadWeb{},
	}

	req := ddosy.ScheduleRequestWeb{
		Endpoint: "test",
		LoadPatterns: []ddosy.LoadPatternWeb{load},
	}

	if err := ddosy.ValidateScheduleRequestWeb(req); err == nil {
		t.Error("expected error got nil")
	}
}

func TestNoLoadPattersScheduleRequestWeb(t *testing.T) {
	traffic := ddosy.TrafficPatternWeb{Weight: 1}

	req := ddosy.ScheduleRequestWeb{
		Endpoint: "test",
		TrafficPatterns: []ddosy.TrafficPatternWeb{traffic},
	}

	if err := ddosy.ValidateScheduleRequestWeb(req); err == nil {
		t.Error("expected error got nil")
	}
}

func TestInvalidLoadPatternScheduleRequestWeb(t *testing.T) {
	load := ddosy.LoadPatternWeb{
		Duration: "-10m",
		Linear: &ddosy.LinearLoadWeb{},
	}
	traffic := ddosy.TrafficPatternWeb{Weight: 1}

	req := ddosy.ScheduleRequestWeb{
		Endpoint: "test",
		LoadPatterns: []ddosy.LoadPatternWeb{load},
		TrafficPatterns: []ddosy.TrafficPatternWeb{traffic},
	}

	if err := ddosy.ValidateScheduleRequestWeb(req); err == nil {
		t.Error("expected error got nil")
	}
}

func TestInvalidTrafficPatternScheduleRequestWeb(t *testing.T) {
	load := ddosy.LoadPatternWeb{
		Duration: "10m",
		Linear: &ddosy.LinearLoadWeb{},
	}
	traffic := ddosy.TrafficPatternWeb{Weight: 0}

	req := ddosy.ScheduleRequestWeb{
		Endpoint: "test",
		LoadPatterns: []ddosy.LoadPatternWeb{load},
		TrafficPatterns: []ddosy.TrafficPatternWeb{traffic},
	}

	if err := ddosy.ValidateScheduleRequestWeb(req); err == nil {
		t.Error("expected error got nil")
	}
}