package ddosy

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"sync"
// 	"sync/atomic"
// 	"time"

// 	vegeta "github.com/tsenart/vegeta/v12/lib"
// )

// func Start(cfg ServerConfig) error {
// 	s := newServer(cfg)
// 	// routing
// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/run", s.runHandler)
// 	mux.HandleFunc("/status", s.statusHandler)
// 	mux.HandleFunc("/kill", s.killHandler)

// 	// server config
// 	srv := &http.Server{
// 		Handler: mux,
// 		Addr:    fmt.Sprintf(":%d", cfg.Port),
// 	}

// 	return srv.ListenAndServe()
// }


// type ServerConfig struct {
// 	Port     int
// 	MaxQueue int
// }

// type server struct {
// 	runningId   int64
// 	tasks       chan *task
// 	runningTask *task

// 	lock *sync.Mutex
// }

// func newServer(cfg ServerConfig) *server {
// 	// instance
// 	s := &server{
// 		tasks: make(chan *task, cfg.MaxQueue),
// 		lock:  &sync.Mutex{},
// 	}

// 	// consume worker
// 	go s.runner()

// 	return s
// }

// func (s *server) runner() {
// 	log.Println("Server started")
// 	for task := range s.tasks {
// 		s.lock.Lock()
// 		s.runningTask = task
// 		s.lock.Unlock()

// 		targeter := NewWeightedTargeter(task.req.TrafficPatterns)
// 		attacker := vegeta.NewAttacker()
// 		main: for _, load := range task.req.LoadPatterns {
// 			pacer, _ := load.Pacer()
// 			for _ = range attacker.Attack(targeter, pacer, load.duration(), "attack") {
// 				if task.shouldStop() {
// 					attacker.Stop()
// 					break main
// 				}
// 			}
// 		}

// 		// TODO metrics

// 		s.lock.Lock()
// 		s.runningTask = nil
// 		s.lock.Unlock()
// 	}
// }

// // add task to a queue
// func (a *server) runHandler(w http.ResponseWriter, r *http.Request) {
// 	var req LoadRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	if err := ValidateLoadRequest(req); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return 
// 	}

// 	resp := a.addToQ(req)
// 	json.NewEncoder(w).Encode(resp)
// }

// func (a *server) addToQ(req LoadRequest) ScheduleResponse {
// 	id := atomic.AddInt64(&a.runningId, 1)
// 	select {
// 	case a.tasks <- &task{id: id, req: req}:
// 		return ScheduleResponse{Id: id}
// 	case <-time.After(time.Millisecond):
// 		return ScheduleResponse{Error: "queue is full"}
// 	}
// }

// // kill currently running task
// func (a *server) killHandler(w http.ResponseWriter, r *http.Request) {
// 	a.lock.Lock()
// 	defer a.lock.Unlock()

// 	t := a.runningTask
// 	if t == nil {
// 		w.Write([]byte("no running tasks"))
// 	} else {
// 		t.stop = true
// 		w.Write([]byte(fmt.Sprintf("killed task with id %d", t.id)))
// 	}
// }

// // what is current status of the task or what are results if the task is done
// func (a *server) statusHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("status"))
// }