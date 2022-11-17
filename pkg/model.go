package ddosy

import (
	"fmt"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type ScheduleRequest struct {
	Endpoint        string           `json:"endpoint"`
	LoadPatterns    []LoadPattern    `json:"load"`
	TrafficPatterns []TrafficPattern `json:"traffic"`
}

func (s ScheduleRequest) Validate() error {
	if s.Endpoint == "" {
		return fmt.Errorf("endpoint is not set")
	}

	if len(s.LoadPatterns) == 0 {
		return fmt.Errorf("load patterns are not set")
	}
	for _, load := range s.LoadPatterns {
		if err := load.Validate(); err != nil {
			return err
		}
	}

	if len(s.TrafficPatterns) == 0 {
		return fmt.Errorf("traffic pattern are not set")
	}
	for _, traffic := range s.TrafficPatterns {
		if err := traffic.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type LoadPattern struct {
	Duration string      `json:"duration"`
	Linear   *LinearLoad `json:"linear"`
	Sine     *SineLoad   `json:"sine"`
}

func (p LoadPattern) Validate() error {
	if t, err := time.ParseDuration(p.Duration); err != nil {
		return err
	} else if t.Seconds() <= 0 {
		return fmt.Errorf("duration must be > 0")
	} else if _, err := p.Pacer(); err != nil {
		return err
	}
	return nil
}

func (p LoadPattern) Pacer() (vegeta.Pacer, error) {
	if p.Linear != nil {
		d, _ := time.ParseDuration(p.Duration)
		return p.Linear.pacer(d), nil
	} else if p.Sine != nil {
		return p.Sine.pacer()
	} else {
		return nil, fmt.Errorf("load pattern must have linear or sine pattern define")
	}
}

type TrafficPattern struct {
	Weight  float64 `json:"weight"`
	Payload string  `json:"payload"`
}

func (p TrafficPattern) Validate() error {
	if p.Weight <= 0 {
		return fmt.Errorf("weight must be > 0")
	}
	return nil
}

type ScheduleResponse struct {
	Id    int64  `json:"id"`
	Error string `json:"error"`
}

type LinearLoad struct {
	StartRate int `json:"startRate"`
	EndRate   int `json:"endRate"`
}

func (l *LinearLoad) pacer(d time.Duration) vegeta.Pacer {
	slope := float64(l.EndRate-l.StartRate) / d.Seconds()
	return vegeta.LinearPacer{
		StartAt: vegeta.Rate{Freq: l.StartRate, Per: time.Second},
		Slope:   slope,
	}
}

type SineLoad struct {
	Mean   int    `json:"mean"`
	Amp    int    `json:"amplitude"`
	Period string `json:"period"`
}

func (s *SineLoad) pacer() (vegeta.Pacer, error) {
	period, err := time.ParseDuration(s.Period)
	if err != nil {
		return nil, err
	}
	if period.Seconds() <= 0 {
		return nil, fmt.Errorf("period must be > 0")
	}
	if s.Mean <= 0 {
		return nil, fmt.Errorf("mean must be > 0")
	}
	if s.Amp >= s.Mean {
		return nil, fmt.Errorf("amplitude must be less then mean")
	}
	return vegeta.SinePacer{
		Period: period,
		Mean:   vegeta.Rate{Freq: s.Mean, Per: time.Second},
		Amp:    vegeta.Rate{Freq: s.Amp, Per: time.Second},
	}, nil

}

type task struct {
	id   int64
	stop bool
}

func (t *task) run() {

}
