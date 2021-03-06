package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const cbatchheader = `
  ________      ________      ________      _________    ________      ___  ___
 |\   ____\    |\   __  \    |\   __  \    |\___   ___\ |\   ____\    |\  \|\  \
 \ \  \___|    \ \  \|\ /_   \ \  \|\  \   \|___ \  \_| \ \  \___|    \ \  \\\  \
  \ \  \        \ \   __  \   \ \   __  \       \ \  \   \ \  \        \ \   __  \
   \ \  \____    \ \  \|\  \   \ \  \ \  \       \ \  \   \ \  \____    \ \  \ \  \
    \ \_______\   \ \_______\   \ \__\ \__\       \ \__\   \ \_______\   \ \__\ \__\
     \|_______|    \|_______|    \|__|\|__|        \|__|    \|_______|    \|__|\|__|

Contain all the jobs!
`

const joboutputheader = `
-----------------------------------------------------------------------------------
| Job Output:                                                                     |
-----------------------------------------------------------------------------------
`

const joboutputend = `
-----------------------------------------------------------------------------------
| End of job output                                                                |
-----------------------------------------------------------------------------------
`

var config Config

func main() {
	fmt.Println(cbatchheader)

	// Measure time from the start
	start := time.Now()

	// Initialize the logger
	initLog()

	// Create the docker driver.
	d, err := NewDockerDriver()
	if err != nil {
		log.Fatal(err)
	}

	// Signal handling
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs)

	go func() {

		for {
			sig := <-sigs

			switch sig {
			case syscall.SIGINT:
				log.Warn("Got signal " + sig.String() + " . Terminating")
				d.Abort()
				os.Exit(0)
			default:
				log.Warn("Got signal " + sig.String() + " . Doing nothing.")
			}
		}
	}()

	// Reads the config
	err = ReadConfig(&config)
	if err != nil {
		log.Fatal(err)
	}

	// Configure the logger
	configureLog()

	// Get job information from environment
	nj := getNewJobFromEnv()

	j, err := nj.CreateJob()
	if err != nil {
		log.Fatal(err)
	}

	err = d.Prepare(j)
	if err != nil {
		log.Fatal(err)
	}

	err = d.Run(j)
	if err != nil {
		log.Fatal(err)
	}

	MeasureTime(start, log.Fields{"job": j, "job_id": j.ID}, "Total time used")

}
