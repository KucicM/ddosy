package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

type App struct {
	runningId int64
	tasks     chan ScheduleRequest
}

func (a *App) Start() {
	go func() {
		log.Println("Server started")
		for task := range a.tasks {
			log.Printf("processing: %d", task.id)
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
	case a.tasks <- ScheduleRequest{
		id: id,
	}:
		return ScheduleResponse{Id: id}
	case <-time.After(time.Second):
		return ScheduleResponse{Error: "queue is full"}
	}
}

func (a *App) KillHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("killed"))
}

func (a *App) StatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("status"))
}
