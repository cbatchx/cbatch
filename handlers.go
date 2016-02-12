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
	log.Printf("New job: %v \n", j)

	//err = j.CreateRunImage()
	// if err != nil {
	//	log.Fatal(err)
	//}

	log.Printf("Started container %v \n", j.Container)
}

func execHandler(w http.ResponseWriter, r *http.Request) {
	j, err := DecodeJob(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	s := fmt.Sprintf("export %v=\"%v\"; export %v=\"%v\"", shellErrorEnvVar, "nil", shellContainerEnvVar, "7123456-erf2-341")
	w.Write([]byte(s))

	log.Printf("Mom is now executing %v \n", j)
}

func doneHandler(w http.ResponseWriter, r *http.Request) {
	j, err := DecodeJob(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Done job: %v \n", j)
}
