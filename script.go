package main

import (
	"errors"
	"os"
)

// Script represents the script of a job
type Script struct {
	path string
	file *os.File
}

// Read
func (s *Script) Read(p []byte) (n int, err error) {
	if s.file == nil {
		return 0, errors.New("Script not opened")
	}

	return s.file.Read(p)
}

// Write
func (s *Script) Write(p []byte) (n int, err error) {
	if s.file == nil {
		return 0, errors.New("Script not opened")
	}

	return s.file.Write(p)
}

// Open the script
func (s *Script) Open() error {
	f, err := os.Open(s.path)
	if err != nil {
		return err
	}

	s.file = f
	return nil
}

// Close the script
func (s *Script) Close() error {
	if s.file == nil {
		return nil
	}

	err := s.file.Close()
	if err != nil {
		return err
	}

	s.file = nil
	return nil
}

// Clear out the content of the file
func (s *Script) Clear() error {
	return ClearFile(s.file)
}
