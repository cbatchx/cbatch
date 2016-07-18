package main

import (
	"time"

	log "github.com/Sirupsen/logrus"
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
const newjobheader = `
-----------------------------------------------------------------------------------
| Job:                                                                            |
-----------------------------------------------------------------------------------
`

const userheader = `
-----------------------------------------------------------------------------------
| User:                                                                           |
-----------------------------------------------------------------------------------
`

const joboutputheader = `
-----------------------------------------------------------------------------------
| Job Output:                                                                     |
-----------------------------------------------------------------------------------
`

var config Config

func main() {
	defer MeasureTime(time.Now(), nil, "Total time used")

	// Reads the config
	err := ReadConfig(&config)
	if err != nil {
		log.Fatal(err)
	}

	// Initlize the logger
	initLog()

	// Get job information from environment
	nj := getNewJobFromEnv()

	j, err := nj.CreateJob()
	if err != nil {
		log.Fatal(err)
	}

	d, err := NewDockerDriver()
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

}
