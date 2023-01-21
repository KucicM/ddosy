package ddosy

import (
	"net/http"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type ScheduleResponseWeb struct {
	Id    uint64 `json:"id"`
	Error string `json:"error"`
}

type ScheduleRequestWeb struct {
	Endpoint        string              `json:"endpoint"`
	LoadPatterns    []LoadPatternWeb    `json:"load"`
	TrafficPatterns []TrafficPatternWeb `json:"traffic"`
}

type LoadPatternWeb struct {
	Duration string         `json:"duration"`
	Linear   *LinearLoadWeb `json:"linear"`
	Sine     *SineLoadWeb   `json:"sine"`
}

type TrafficPatternWeb struct {
	Weight  float64 `json:"weight"`
	Payload string  `json:"payload"`
}

type LinearLoadWeb struct {
	StartRate int `json:"startRate"`
	EndRate   int `json:"endRate"`
}

type SineLoadWeb struct {
	Mean   int    `json:"mean"`
	Amp    int    `json:"amplitude"`
	Period string `json:"period"`
}

// database
type DatabaseTask struct {
	Id        uint64
	StatusId  TaskStatus
	CreatedAt time.Time
	StartedAt *time.Time
	KilledAt  *time.Time
	DoneAt    *time.Time
	Request   ScheduleRequestWeb
	Results   string
}

type TaskStatus int8

const (
	Scheduled TaskStatus = iota + 1
	Running
	Killed
	Done
)

// internal
type LoadTask struct {
	id      uint64
	traffic TrafficPattern
	load    []LoadPattern
}

type LoadPattern struct {
	duration time.Duration
	pacer    vegeta.Pacer
}

type TrafficPattern struct {
	endpoint string
	header   http.Header
	method   string
	dist     TrafficDistribution
}

type TrafficDistribution struct {
	weigths  []float64
	payloads [][]byte
}

func NewLoadTask(req ScheduleRequestWeb) *LoadTask {
	load := make([]LoadPattern, len(req.LoadPatterns))
	for i, p := range req.LoadPatterns {
		load[i] = NewLoadPattern(p)
	}

	return &LoadTask{
		traffic: NewTrafficPattern(req.Endpoint, req.TrafficPatterns),
		load:    load,
	}
}

func NewTrafficPattern(endpoint string, patterns []TrafficPatternWeb) TrafficPattern {

	ws := make([]float64, len(patterns))
	payloads := make([][]byte, len(patterns))
	for i, p := range patterns {
		ws[i] = p.Weight
		payloads[i] = []byte(p.Payload)
	}

	dist := TrafficDistribution{
		weigths:  CreateDistribution(ws),
		payloads: payloads,
	}

	header := http.Header{}
	header.Add("Content-Type", "application/json")
	return TrafficPattern{
		endpoint: endpoint,
		header:   header,
		method:   http.MethodPost,
		dist:     dist,
	}
}

func NewLoadPattern(pattern LoadPatternWeb) LoadPattern {
	duration, _ := time.ParseDuration(pattern.Duration)
	var pacer vegeta.Pacer

	if pattern.Linear != nil {
		l := pattern.Linear
		slope := float64(l.EndRate-l.StartRate) / duration.Seconds()
		pacer = vegeta.LinearPacer{
			StartAt: vegeta.Rate{Freq: l.StartRate, Per: time.Second},
			Slope:   slope,
		}
	} else {
		s := pattern.Sine
		p, _ := time.ParseDuration(s.Period)
		pacer = vegeta.SinePacer{
			Period: p,
			Mean:   vegeta.Rate{Freq: s.Mean, Per: time.Second},
			Amp:    vegeta.Rate{Freq: s.Amp, Per: time.Second},
		}
	}

	return LoadPattern{
		duration: duration,
		pacer:    pacer,
	}
}

func (t *LoadTask) Targeter() vegeta.Targeter {
	return NewWeightedTargeter(t.traffic)
}
