package main

type ScheduleRequest struct {
	id int64
}

type ScheduleResponse struct {
	Id int64 `json:"id"`
	Error string `json:"error"`
}