package ddosy

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"
	"sync"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type ResultProvider struct {
	repo           *TaskRepository
	lock           *sync.RWMutex
	runningMetrics map[uint64]*vegeta.Metrics
}

func NewRelustProvider(repo *TaskRepository) *ResultProvider {
	p := &ResultProvider{
		repo:           repo,
		lock:           &sync.RWMutex{},
		runningMetrics: make(map[uint64]*vegeta.Metrics),
	}
	return p
}

// updates in memory metrics
func (p *ResultProvider) UpdateRunning(id uint64, res *vegeta.Result) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if r, ok := p.runningMetrics[id]; ok && r != nil {
		r.Add(res)
	} else {
		m := &vegeta.Metrics{}
		m.Add(res)
		p.runningMetrics[id] = m
	}
}

// flush in memory metris to database
func (p *ResultProvider) FinalizeRunning(id uint64) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if m, ok := p.runningMetrics[id]; ok && m != nil {
		str := metricsToStr(m)
		p.repo.UpdateProgress(id, str)
		delete(p.runningMetrics, id)
	}
}

func (p *ResultProvider) Get(id uint64) (string, error) {
	task, err := p.repo.Get(id)
	if err != nil {
		log.Printf("error getting task with id=%d %s\n", id, err)
		return "", err
	}

	p.lock.RLock()
	var current string
	if m, ok := p.runningMetrics[id]; ok {
		current = metricsToStr(m)
	}
	p.lock.RUnlock()

	res := LoadTaskReport{}
	res.Endpoint = task.Request.Endpoint
	res.Metrics = strings.Join([]string{task.Results, current}, "\n")
	res.Status = TaskStatusStr[task.StatusId]
	res.CreatedAt = task.CreatedAt.String()
	if task.StartedAt != nil {
		res.StartedAt = task.StartedAt.String()
	}
	if task.KilledAt != nil {
		res.KilledAt = task.KilledAt.String()
	}
	if task.DoneAt != nil {
		res.DoneAt = task.DoneAt.String()
	}
	if b, err := json.MarshalIndent(task.Request, "", " "); err == nil {
		res.Request = string(b)
	}
	return res.String()
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
