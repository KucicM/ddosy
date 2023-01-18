package ddosy

type ScheduleResponseWeb struct {
	Id    int64  `json:"id"`
	Error string `json:"error"`
}

type ScheduleRequestWeb struct {
	Endpoint        string           `json:"endpoint"`
	LoadPatterns    []LoadPatternWeb    `json:"load"`
	TrafficPatterns []TrafficPatternWeb `json:"traffic"`
}

type LoadPatternWeb struct {
	Duration string      `json:"duration"`
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

// func (p LoadPattern) Pacer() (vegeta.Pacer, error) {
// 	if p.Linear != nil {
// 		d := p.duration()
// 		return p.Linear.pacer(d), nil
// 	} else if p.Sine != nil {
// 		return p.Sine.pacer()
// 	} else {
// 		return nil, fmt.Errorf("load pattern must have linear or sine pattern define")
// 	}
// }

// func (p LoadPattern) duration() time.Duration {
// 	d, _ := time.ParseDuration(p.Duration)
// 	return d
// }


// func (l *LinearLoad) pacer(d time.Duration) vegeta.Pacer {
// 	slope := float64(l.EndRate-l.StartRate) / d.Seconds()
// 	return vegeta.LinearPacer{
// 		StartAt: vegeta.Rate{Freq: l.StartRate, Per: time.Second},
// 		Slope:   slope,
// 	}
// }


// func (s *SineLoad) pacer() (vegeta.Pacer, error) {
// 	period, err := time.ParseDuration(s.Period)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if period.Seconds() <= 0 {
// 		return nil, fmt.Errorf("period must be > 0")
// 	}
// 	if s.Mean <= 0 {
// 		return nil, fmt.Errorf("mean must be > 0")
// 	}
// 	if s.Amp >= s.Mean {
// 		return nil, fmt.Errorf("amplitude must be less then mean")
// 	}
// 	return vegeta.SinePacer{
// 		Period: period,
// 		Mean:   vegeta.Rate{Freq: s.Mean, Per: time.Second},
// 		Amp:    vegeta.Rate{Freq: s.Amp, Per: time.Second},
// 	}, nil

// }

// // type task struct {
// // 	id   int64
// // 	stop bool
// // 	req LoadRequest
// // 	stopLock *sync.RWMutex
// // }

// // func (t *task) shouldStop() bool {
// // 	t.stopLock.RLock()
// // 	defer t.stopLock.RUnlock()
// // 	return t.stop
// // }

// // func (t *task) forceStop() {
// // 	t.stopLock.Lock()
// // 	t.stop = true
// // 	t.stopLock.Unlock()
// // }
