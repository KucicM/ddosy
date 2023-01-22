package main

import (
	"flag"
	"log"

	ddosy "github.com/kucicm/ddosy/app"
)

func main() {
	port := flag.Int("port", 4000, "port for the server")
	dbUrl := flag.String("dbUrl", "test.db", "connetion to database")
	flag.Parse()

	cfg := ddosy.ServerConfig{
		Port:              *port,
		DbUrl:             *dbUrl,
		TruncateDbOnStart: false,
	}

	log.Fatalln(ddosy.Start(cfg))
}
