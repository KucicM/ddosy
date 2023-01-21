package ddosy

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"sync"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type ResultProvider struct {
	lock    *sync.RWMutex
	metrics map[uint64]*RunningResult
}

type RunningResult struct {
	current *vegeta.Metrics
	total   []string
}

func NewRelustProvider() *ResultProvider {
	p := &ResultProvider{
		metrics: make(map[uint64]*RunningResult),
		lock:    &sync.RWMutex{},
	}
	// todo persister
	return p
}

func (p *ResultProvider) NewPattern(id uint64) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if m, ok := p.metrics[id]; ok {
		m.total = append(m.total, metricsToStr(m.current))
		m.current = &vegeta.Metrics{}
		return
	}

	p.metrics[id] = &RunningResult{
		current: &vegeta.Metrics{},
		total:   make([]string, 0),
	}
}

func (p *ResultProvider) Update(id uint64, res *vegeta.Result) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if r, ok := p.metrics[id]; ok {
		r.current.Add(res)
	} else {
		log.Printf("cannot find results with id %d\n", id)
	}
}

func (p *ResultProvider) Done(id uint64) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if m, ok := p.metrics[id]; ok {
		m.total = append(m.total, metricsToStr(m.current))
		m.current = nil
	}
}

func metricsToStr(m *vegeta.Metrics) string {
	if m == nil {
		return ""
	}

	buf := bytes.Buffer{}
	m.Close()
	rep := vegeta.NewTextReporter(m)
	if err := rep.Report(&buf); err != nil {
		log.Printf("error writing report to buffer %s\n", err)
		return err.Error()
	}
	return buf.String()
}

func (p *ResultProvider) Get(id uint64) (string, error) {
	p.lock.RLock()
	defer p.lock.RUnlock()
	if m, ok := p.metrics[id]; ok {
		return strings.Join(m.total, "\n\n"), nil
	} else {
		return "", fmt.Errorf("cannot find record with id=%d", id)
	}
}
