package main

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
)

const (
	momlock          = "mom.lock"
	configname       = "config"
	configbackupname = "config.bak"
	jobStartername   = "jobstarter"
	prologueName     = "prologue"
	epilogueName     = "epilogue"
)

var oldConfig []byte

func placeJobStarterScript() error {
	err := writeScript(mompriv+jobStartername,
		jobstartersh,
		0555)
	if err != nil {
		return err
	}

	log.Println("Created Jobstarter file: " + mompriv + jobStartername)
	return nil
}

func removeJobStarterScript() error {
	err := os.Remove(mompriv + jobStartername)
	if err != nil {
		return err
	}

	log.Println("Deleted Jobstarter file: " + mompriv + jobStartername)
	return nil
}

func getMomPid() (int, error) {
	b, err := ioutil.ReadFile(mompriv + momlock)
	if err != nil {
		return -1, err
	}

	s := string(b)
	s = strings.Split(s, "\n")[0]

	pid, err := strconv.Atoi(s)
	if err != nil {
		return -1, err
	}

	return pid, err
}

func momReloadConfig() error {
	pid, err := getMomPid()
	if err != nil {
		log.Println("Could not retrive PID of MOM daemon is it running?")
		return err
	}

	err = syscall.Kill(pid, syscall.SIGHUP)
	if err != nil {
		return err
	}

	log.Printf("Sent SIGHUP to %v to make mom reload config.\n", pid)
	return nil
}

func configChange() error {
	b, err := ioutil.ReadFile(mompriv + configname)
	if err != nil {
		return err
	}

	oldConfig = b
	newConfig := string(oldConfig) + "$job_starter " + mompriv + jobStartername + "\n"

	err = ioutil.WriteFile(mompriv+configname, []byte(newConfig), 0644)
	if err != nil {
		return err
	}
	log.Printf("Succesfully changed configfile %v \n", mompriv+configname)

	err = momReloadConfig()
	if err != nil {
		return err
	}

	return nil
}

func configRevert() error {
	err := ioutil.WriteFile(mompriv+configname, oldConfig, 0644)
	if err != nil {
		return err
	}

	err = momReloadConfig()
	if err != nil {
		return err
	}

	log.Printf("Succesfully reverted configfile %v \n", mompriv+configname)
	return nil
}

func placePrologueScript() error {
	err := writeScript(mompriv+prologueName,
		prologuesh,
		0500)
	if err != nil {
		return err
	}

	log.Println("Created Prologue file: " + mompriv + prologueName)
	return nil
}

func removePrologueScript() error {
	err := os.Remove(mompriv + prologueName)
	if err != nil {
		return err
	}

	log.Println("Deleted Prologue file: " + mompriv + prologueName)
	return nil
}

func placeEpilogueScript() error {
	err := writeScript(mompriv+epilogueName, epiloguesh, 0500)
	if err != nil {
		return err
	}

	log.Println("Created Epilogue file: " + mompriv + epilogueName)
	return nil
}

func removeEpilogueScript() error {
	err := os.Remove(mompriv + epilogueName)
	if err != nil {
		return err
	}

	log.Println("Deleted Epilogue file: " + mompriv + epilogueName)
	return nil
}

func writeScript(path, script string, perm os.FileMode) error {
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		return err
	}

	err = f.Chmod(perm)
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
// If anything fails we just quit right away, as the docquer cannot work without
// these "reporters". They are just normal shell scripts that does http requests
// to the built-in http server.
func PlaceReporters() {
	err := placePrologueScript()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	err = placeJobStarterScript()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	err = placeEpilogueScript()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	err = configChange()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

// RemoveReporters removes the prologue and epilogue scripts
// We want to clean up as much as possible so no os.Exit() before end here.
func RemoveReporters() {
	errOccured := false
	err := removePrologueScript()
	if err != nil {
		log.Fatal(err)
		errOccured = true
	}
	err = removeJobStarterScript()
	if err != nil {
		log.Fatal(err)
		errOccured = true
	}
	err = removeEpilogueScript()
	if err != nil {
		log.Fatal(err)
		errOccured = true
	}
	err = configRevert()
	if err != nil {
		log.Fatal(err)
		errOccured = true
	}

	if errOccured {
		os.Exit(1)
	}
}
