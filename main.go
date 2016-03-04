package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
)

const socketpath = "/var/run/cbatch.sock"

func cleanUp() {
	log.Println("Got signal cleaning up!")
	os.Remove(socketpath)
	// RemoveReporters()
}

func main() {
	// PlaceReporters()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, os.Kill)

	go func() {
		<-sigchan
		cleanUp()
		os.Exit(0)
	}()

	js := NewJobStore()
	http.HandleFunc("/new", newHandler(js))
	http.HandleFunc("/exec", execHandler(js))
	http.HandleFunc("/done", doneHandler(js))

	// Clean old socket
	os.Remove(socketpath)
	l, err := net.Listen("unix", "/var/run/cbatch.sock")
	if err != nil {
		log.Fatal(err)
	}
	http.Serve(l, nil)
}
