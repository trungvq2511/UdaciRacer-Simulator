package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

var dataFile string
var addr string

func init() {
	dataFile = os.Getenv("DATA_FILE")
	if dataFile == "" {
		dataFile = "data.json"
	}

	addr = os.Getenv("ADDR")
	if addr == "" {
		addr = "0.0.0.0:3001"
	}
}

func main() {
	if err := generateTrackSegments(dataFile); err != nil {
		log.Fatal(err)
	}

	service, err := NewRaceService(DataFromJSONFile(dataFile))
	panicErr(err)

	server, err := NewServer(service)
	panicErr(err)

	srv := &http.Server{
		Handler:      server,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	srv.ListenAndServe()
}
