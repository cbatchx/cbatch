package main

import (
	"html/template"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
)

// Bootstrap describes a bootstrap script.
type Bootstrap struct {
	tmpfile string
}

const templateFile = "config/bootstrap.tmpl.sh"

// NewBootstrap generates a boostrap script from the boostrap template.
// After the boostrap struct is initalized it can be passed to for instance a driver.
func NewBootstrap(j *Job) (*Bootstrap, error) {

	// Read template.
	var t = template.Must(template.ParseFiles(templateFile))

	// Create temporary file.
	tmpfile, err := ioutil.TempFile("", "cbatch_prepare")
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// Write template to temporary file.
	t.Execute(tmpfile, j)

	// Release the tmpfile
	if err := tmpfile.Close(); err != nil {
		log.Error(err)
		return nil, err
	}

	return &Bootstrap{
		tmpfile: tmpfile.Name(),
	}, nil
}

// GetScriptPath get the path to the temporary script.
func (b *Bootstrap) GetScriptPath() string {
	return b.tmpfile
}

// Remove the temporary bootstrap file.
func (b *Bootstrap) Remove() error {
	return os.Remove(b.tmpfile)
}
