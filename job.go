package main

import "fmt"

// Mount represents a mount point inside a container
type Mount struct {
	Source      string
	Destination string
	RW          bool
}

// Mounts array of Mount
type Mounts []Mount

// Job A runnable job.
type Job struct {
	ID     string   // ID of the job.
	Cmd    []string // Cmd to run on start.
	User   *User    // The user running the job.
	Shell  *Shell   // The shell for the job.
	Image  *Image   // The Image for the job.
	Env    []string // Environment variables on the form KEY=VALUE.
	Mounts Mounts   // Mounts needed for the job to run.
}

// AddMount to mounts
func (m Mounts) AddMount(source, destination string, rw bool) Mounts {
	return append(m, Mount{Source: source, Destination: destination, RW: rw})
}

// DockerBindString returns an docker string representation of a mount
// On the form:
// $HOST:$CONTAINER:ro
func (m Mount) DockerBindString() string {
	if m.RW {
		return fmt.Sprintf("%v:%v", m.Source, m.Destination)
	}
	return fmt.Sprintf("%v:%v:ro", m.Source, m.Destination)
}
