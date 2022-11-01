package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type App struct {
	runningId   int64
	tasks       chan *task
	runningTask *task

	lock *sync.Mutex
}

func (a *App) Start() {
	go func() {
		log.Println("Server started")
		for task := range a.tasks {

			a.lock.Lock()
			a.runningTask = task
			a.lock.Unlock()

			task.run()

			a.lock.Lock()
			a.runningTask = nil
			a.lock.Unlock()
		}
	}()
}

func (a *App) ScheduelHandler(w http.ResponseWriter, r *http.Request) {
	resp := a.schedule()
	json.NewEncoder(w).Encode(resp)
}

func (a *App) schedule() ScheduleResponse {
	id := atomic.AddInt64(&a.runningId, 1)
	select {
	case a.tasks <- &task{
		id: id,
	}:
		return ScheduleResponse{Id: id}
	case <-time.After(time.Second):
		return ScheduleResponse{Error: "queue is full"}
	}
}

func (a *App) KillHandler(w http.ResponseWriter, r *http.Request) {
	a.lock.Lock()
	defer a.lock.Unlock()

	t := a.runningTask
	if t == nil {
		w.Write([]byte("no running tasks"))
	} else {
		t.stop = true
		w.Write([]byte(fmt.Sprintf("killed task with id %d", t.id)))
	}
}

func (a *App) StatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("status"))
}
