package main

import (
	"flag"
	"log"

	ddosy "github.com/kucicm/ddosy/pkg"
)

func main() {
	port := flag.Int("port", 4000, "port for the server")
	queue := flag.Int("queue", 1, "max queue size")
	flag.Parse()

	cfg := ddosy.ServerConfig{
		Port:     *port,
		MaxQueue: *queue,
	}

	// app := ddosy.NewServer(cfg)
	log.Fatalln(ddosy.Start(cfg))
	// log.Fatalln(app.ListenAndServe())

}
