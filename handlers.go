package main

import (
	"fmt"
	"log"
	"net/http"
)

func newHandler(w http.ResponseWriter, r *http.Request) {
	j, err := DecodeJob(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[JOB: %v] New job: %v \n", j.ID, j)

	// TODO move to config
	j.ImageName = "centos:latest"

	err = j.CreateImage()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[JOB: %v] Built image for job: %v \n", j.ID, j)

	err = j.StartContainer()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[JOB: %v] Started container: %v \n", j.ID, j.Container)
}

func execHandler(w http.ResponseWriter, r *http.Request) {
	j, err := DecodeJob(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	s := fmt.Sprintf("export %v=\"%v\"; export %v=\"%v\"", shellErrorEnvVar, "nil", shellContainerEnvVar, "7123456-erf2-341")
	w.Write([]byte(s))

	log.Printf("[JOB: %v] Mom is now executing %v \n", j.ID, j)
}

func doneHandler(w http.ResponseWriter, r *http.Request) {
	j, err := DecodeJob(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[JOB: %v] Done job: %v \n", j.ID, j)
}
