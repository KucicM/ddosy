package ddosy

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
)

type ServerConfig struct {
	Port     int
	MaxQueue int
}

type Server struct {
	runningId   int64
	tasks       chan *task
	runningTask *task

	lock *sync.Mutex
	srv  *http.Server
}

func NewServer(cfg ServerConfig) *Server {

	// instance
	s := &Server{
		tasks: make(chan *task, cfg.MaxQueue),
		lock:  &sync.Mutex{},
	}

	// routing
	r := mux.NewRouter()
	r.PathPrefix("/schedule").HandlerFunc(s.scheduelHandler).Methods("POST")
	r.PathPrefix("/status").HandlerFunc(s.statusHandler).Methods("GET")
	r.PathPrefix("/kill").HandlerFunc(s.killHandler).Methods("DEL")

	// server config
	s.srv = &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf(":%d", cfg.Port),
	}

	// consume worker
	go func() {
		log.Println("Server started")
		for task := range s.tasks {

			s.lock.Lock()
			s.runningTask = task
			s.lock.Unlock()

			task.run()

			s.lock.Lock()
			s.runningTask = nil
			s.lock.Unlock()
		}
	}()

	return s
}

func (s *Server) ListenAndServe() error {
	return s.srv.ListenAndServe()
}

// add task to a queue
func (a *Server) scheduelHandler(w http.ResponseWriter, r *http.Request) {
	var req ScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := a.schedule(req)
	json.NewEncoder(w).Encode(resp)
}

func (a *Server) schedule(req ScheduleRequest) ScheduleResponse {
	if err := req.Validate(); err != nil {
		return ScheduleResponse{Error: err.Error()}
	}

	id := atomic.AddInt64(&a.runningId, 1)
	select {
	case a.tasks <- &task{id: id}:
		return ScheduleResponse{Id: id}
	case <-time.After(time.Second):
		return ScheduleResponse{Error: "queue is full"}
	}
}

// kill currently running task
func (a *Server) killHandler(w http.ResponseWriter, r *http.Request) {
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

// what is current status of the task or what are results if the task is done
func (a *Server) statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("status"))
}
