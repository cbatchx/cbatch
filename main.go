package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
)

func cleanUp() {
	log.Println("Got signal cleaning up!")
	RemoveReporters()
}

func main() {
	PlaceReporters()

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
	http.ListenAndServe(":8080", nil)
}
