package main

import "net/http"

type App struct {
}

func (a *App) ScheduelHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("scheduled"))
}

func (a *App) KillHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("killed"))
}

func (a *App) StatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("status"))
}
