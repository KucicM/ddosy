package ddosy

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"net/http"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type ServerConfig struct {
	Port              int
	DbUrl             string
	TruncateDbOnStart bool
}

func Start(cfg ServerConfig) error {
	srv := NewServer(cfg)

	mux := http.NewServeMux()
	mux.HandleFunc("/run", srv.ScheduleHandler)
	mux.HandleFunc("/status", srv.StatusHandler)
	mux.HandleFunc("/kill", srv.KillHandler)

	// server config
	httpSrv := &http.Server{
		Handler: mux,
		Addr:    fmt.Sprintf(":%d", cfg.Port),
	}

	return httpSrv.ListenAndServe()
}

type Server struct {
	taskProvider   *TaskProvider
	resultProvider *ResultProvider
	kill           chan struct{}
}

func NewServer(cfg ServerConfig) *Server {
	repo := NewTaskRepository(cfg.DbUrl, cfg.TruncateDbOnStart) // todo close connection on shutdown

	// todo shutdown

	srv := &Server{
		taskProvider:   NewTaskProvider(repo),
		resultProvider: NewRelustProvider(repo),
		kill:           make(chan struct{}, 1),
	}

	go srv.runner()

	return srv
}

func (s *Server) runner() {
	log.Println("server runner started")
	for {
		task := s.taskProvider.Next()
		if task == nil {
			time.Sleep(time.Second)
			continue
		}

		log.Printf("Running task %+v", task)
		targeter := task.Targeter()
		attacker := vegeta.NewAttacker()

	main:
		for _, load := range task.load {
			for res := range attacker.Attack(targeter, load.pacer, load.duration, "attack") {
				select {
				case <-s.kill:
					attacker.Stop()
					log.Printf("load test with id=%d killed\n", task.id)
					s.taskProvider.Kill(task.id)
					break main
				default:
					s.resultProvider.Update(task.id, res)
				}
			}
		}
		s.taskProvider.Done(task.id)
	}
}

func (s *Server) ScheduleHandler(w http.ResponseWriter, r *http.Request) {
	var req ScheduleRequestWeb
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := ValidateScheduleRequestWeb(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if id, err := s.taskProvider.ScheduleTask(req); err == nil {
		json.NewEncoder(w).Encode(ScheduleResponseWeb{Id: id})
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) KillHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("recived kill request")
	select {
	case s.kill <- struct{}{}:
		log.Println("kill signal sent")
		w.Write([]byte("Stopping"))
	case <-time.After(time.Second):
		w.Write([]byte("No running tasks?"))
	}
}

func (s *Server) StatusHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := s.resultProvider.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(result))
}
