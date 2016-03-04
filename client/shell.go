package main

import "io"

// Shell represents a shell
type Shell struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader
	TTY    bool // Is the job Interactive.
}
