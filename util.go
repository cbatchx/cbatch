package main

import "os"

// ClearFile Remove all contents of a file
func ClearFile(f *os.File) error {
	err := f.Truncate(0)
	if err != nil {
		return err
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		return err
	}

	return nil
}
