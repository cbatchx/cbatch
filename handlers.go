package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func newHandler(w http.ResponseWriter, r *http.Request) {
	j, err := DecodeJob(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	s := j.GetScript()
	s.Open()
	defer s.Close()
	b, err := ioutil.ReadAll(s)

	log.Printf("New job: %v \n", j)
	log.Printf("Shell script:\n%v", string(b))
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
