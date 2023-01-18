package ddosy

import (
	"fmt"
	"sync/atomic"
	"time"
)

type TaskProvider struct {
	runningId uint64
	tasks chan *LoadTask
}

func NewTaskProvider(queueSize int) *TaskProvider {
	return &TaskProvider{
		tasks: make(chan *LoadTask, queueSize),
	}
}

func (p *TaskProvider) ScheduleTask(task *LoadTask) (uint64, error) {
	task.id = atomic.AddUint64(&p.runningId, 1)
	select {
	case p.tasks <- task:
		return task.id, nil
	case <- time.After(time.Second):
		return 0, fmt.Errorf("queue is full")
	}
}
