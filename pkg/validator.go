package ddosy

import (
	"fmt"
	"time"
)

func ValidateScheduleRequestWeb(req ScheduleRequestWeb) error {
	if req.Endpoint == "" {
		return fmt.Errorf("must provide endpoint")
	}

	if len(req.TrafficPatterns) == 0 {
		return fmt.Errorf("must provide at least 1 traffic pattern")
	}

	for _, p := range req.TrafficPatterns {
		if err := ValidateTrafficPatternWeb(p); err != nil {
			return err
		}
	}

	if len(req.LoadPatterns) == 0 {
		return fmt.Errorf("must provide at least 1 load pattern")
	}

	for _, p := range req.LoadPatterns {
		if err := ValidateLoadPatternWeb(p); err != nil {
			return err
		}
	}

	return nil
}

func ValidateTrafficPatternWeb(pattern TrafficPatternWeb) error {
	if pattern.Weight <= 0 {
		return fmt.Errorf("invalid traffic weight")
	}
	return nil
}


func ValidateLoadPatternWeb(pattern LoadPatternWeb) error {
	duration, err := time.ParseDuration(pattern.Duration)
	if err != nil {
		return err
	}

	if duration.Seconds() < 0 {
		return fmt.Errorf("load duration must be > 0s")
	}

	if pattern.Linear != nil {
		return ValidateLinearLoadWeb(*pattern.Linear)
	}

	if pattern.Sine != nil {
		return ValidateSineLoadWeb(*pattern.Sine)
	}

	return fmt.Errorf("must define linear or sine pattern")
}

func ValidateLinearLoadWeb(load LinearLoadWeb) error {
	if load.StartRate < 0 {
		return fmt.Errorf("start rate cannot be less then 0")
	}

	if load.EndRate < 0 {
		return fmt.Errorf("end rate cannot be less then 0")
	}

	return nil
}

func ValidateSineLoadWeb(load SineLoadWeb) error {
	period, err := time.ParseDuration(load.Period)
	if err != nil {
		return err
	}

	if period.Seconds() <= 0 {
		return fmt.Errorf("period must be > 0")
	}

	if load.Mean <= 0 {
		return fmt.Errorf("mean must be > 0")
	}

	if load.Amp < 0 {
		return fmt.Errorf("amplitude must be >= 0")
	}

	if load.Amp >= load.Mean {
		return fmt.Errorf("amplitude must be less then mean")
	}

	return nil
}
