package main

// Driver Interface for container drivers
type Driver interface {
	Run(j *Job) error
}
