package main

import (
	"fmt"
	"log"
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

	// Reads the config
	err := ReadConfig(&config)
	if err != nil {
		log.Fatal(err)
	}

	nj := getNewJobFromEnv()
	fmt.Printf(cbatchheader)

	fmt.Printf(newjobheader)
	fmt.Printf("%+v \n", nj)

	j, err := nj.CreateJob()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(userheader)
	fmt.Printf("%+v\n", j.User)

	fmt.Printf(joboutputheader)

	d := NewDockerDriver()

	err = d.Run(j)
	if err != nil {
		log.Fatal(err)
	}

}
