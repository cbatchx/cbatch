package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"bitbucket.org/dizk/cbatch/types"
)

func newHandler(js *JobStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var nj types.NewJob
		d := json.NewDecoder(r.Body)
		err := d.Decode(&nj)
		if err != nil {
			log.Fatal(err)
		}

		j, err := DecodeJob(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("[JOB: %v] New job: %v \n", j.ID, j)

		j.InitJob()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("[JOB: %v] Got info about the user: %v \n", j.ID, j.User)

		// TODO move to config
		j.ImageName = "centos:latest"

		err = j.CreateImage()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("[JOB: %v] Built image for job: %v \n", j.ID, j)

		// Save the job
		js.Save(j)

	}
}

func execHandler(js *JobStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpj, err := DecodeJob(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		// Get the job
		j := js.Get(tmpj.ID)

		if j.ID == "" {
			log.Printf("[JOB: %v] No job found \n", "")
			return
		}

		j.Cmd = tmpj.Cmd
		log.Printf("[JOB: %v] Starting container: %v \n", j.ID, j.Container)
		err = j.StartContainer()
		if err != nil {
			log.Fatal(err)
		}

		js.Save(j)

		s := fmt.Sprintf("export %v=\"%v\"; export %v=\"%v\"", shellErrorEnvVar, "nil", shellContainerEnvVar, j.Container.ID)
		w.Write([]byte(s))

		log.Printf("[JOB: %v] Executing job: %v \n", j.ID, j)
	}
}

func doneHandler(js *JobStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpj, err := DecodeJob(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		j := js.Get(tmpj.ID)
		log.Printf("[JOB: %v] Done job: %v \n", j.ID, j)

		log.Printf("[JOB: %v] Stopping container: %v \n", j.ID, j.Container)
		err = j.StopContainer()
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("[JOB: %v] Removing container: %v \n", j.ID, j.Container)
		err = j.RemoveContainer()
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("[JOB: %v] Removing image: %v \n", j.ID, j.ImageName)
		err = j.RemoveImage()
		if err != nil {
			log.Fatal(err)
		}
	}
}
