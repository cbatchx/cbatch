package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func newJobHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var j Job
	err = json.Unmarshal(b, &j)
	if err != nil {
		log.Fatal(err)
	}

	s := j.GetScript()
	s.Open()
	defer s.Close()

	b, err = ioutil.ReadAll(s)

	log.Printf("New job: %v \n", j)
	log.Printf("Shell script %v", string(b))
}

func doneJobHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var j Job
	err = json.Unmarshal(b, &j)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Done job: %v \n", j)
}

func cleanUp() {
	RemoveReporters()
}

func main() {
	PlaceReporters()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, os.Kill)

	http.HandleFunc("/new", newJobHandler)
	http.HandleFunc("/done", doneJobHandler)
	http.ListenAndServe(":8080", nil)

	go func() {
		<-sigchan
		// When we get a signal just exit and Clean up

		cleanUp()
		os.Exit(0)
	}()
}
