package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func newHandler(w http.ResponseWriter, r *http.Request) {
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

func execHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Mom is now executing %v \n", string(b))
}

func doneHandler(w http.ResponseWriter, r *http.Request) {
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
