package main

type ScheduleRequest struct {
}

type ScheduleResponse struct {
	Id int64 `json:"id"`
	Error string `json:"error"`
}

type task struct {
	id int64
	stop bool
}

func (t *task) run() {

}