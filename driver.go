package main

// Driver Interface for container drivers
type Driver interface {
	// Prepare pulls the image and sets a new image to use on the job.
	Prepare(j *Job) error
	// Runs the prepared image.
	Run(j *Job) error
	// Abort
	Abort(j *Job) error
}
