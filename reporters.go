package main

import (
	"log"
	"os"
)

const (
	mompriv    = "/var/spool/torque/mom_priv/"
	curl       = `curl -H "Content-Type: application/json" -X POST -d`
	newPath    = "http://localhost:8080/new \n"
	donePath   = "http://localhost:8080/done \n"
	bashHeader = "#!/bin/sh \n"
	// http://docs.adaptivecomputing.com/torque/6-0-0/help.htm#topics/torque/13-appendices/scriptEnvironment.htm
	prologueJSON = ` "{ \"jobid\":\"$1\", \"user\":\"$2\", \"group\":\"$3\", \"jobname\":\"$4\", \"resourcelimits\":\"$5\", \"jobqueue\":\"$6\", \"jobaccount\":\"$7\" }" `
	prologueName = "prologue"
	epilogueJSON = ` "{ \"jobid\":\"$1\", \"user\":\"$2\", \"group\":\"$3\", \"jobname\":\"$4\", \"sessionid\":\"$5\", \"resourcelimits\":\"$6\", \"resourcesused\":\"$7\", \"jobqueue\":\"$8\", \"jobaccount\":\"$9\", \"jobexitcode\":\"$10\" }" `
	epilogueName = "epilogue"
)

func placePrologueScript() {
	err := writeScript(mompriv+prologueName, bashHeader+curl+prologueJSON+newPath)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Created Prologue file: " + mompriv + prologueName)
}

func removePrologueScript() {
	err := os.Remove(mompriv + prologueName)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Deleted Prologue file: " + mompriv + prologueName)
}

func placeEpilogueScript() {
	err := writeScript(mompriv+epilogueName, bashHeader+curl+epilogueJSON+donePath)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Created Epilogue file: " + mompriv + epilogueName)
}

func removeEpilogueScript() {
	err := os.Remove(mompriv + epilogueName)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Deleted Epilogue file: " + mompriv + epilogueName)
}

func writeScript(path, script string) error {
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		return err
	}

	err = f.Chmod(0500)
	if err != nil {
		return err
	}

	_, err = f.WriteString(script)
	if err != nil {
		return err
	}

	return nil
}

// PlaceReporters places the Prologue and Epilogue scripts
func PlaceReporters() {
	placePrologueScript()
	placeEpilogueScript()
}

// RemoveReporters removes the prologue and epilogue scripts
func RemoveReporters() {
	removePrologueScript()
	removeEpilogueScript()
}
