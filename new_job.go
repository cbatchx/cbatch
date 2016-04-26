package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

// NewJob represents a new job
type NewJob struct {
	PBSJob PBSJob
	Args   []string
}

// CreateJob creates a new Job.
func (n *NewJob) CreateJob() (*Job, error) {

	// Get the user
	u, err := n.GetUser()
	if err != nil {
		return nil, err
	}

	// Get the image
	i, err := n.GetImage()
	if err != nil {
		return nil, err
	}

	// Get the shell
	s, err := n.GetShell()
	if err != nil {
		return nil, err
	}

	var m Mounts
	m, err = n.addMounts(m)
	if err != nil {
		return nil, err
	}

	e, err := n.GetEnv()
	if err != nil {
		return nil, err
	}

	return &Job{
		User:   u,
		Cmd:    n.Args,
		Shell:  s,
		Image:  i,
		Mounts: m,
		Env:    e,
	}, nil
}

// GetEnv get the current Environment
func (n *NewJob) GetEnv() ([]string, error) {
	e := os.Environ()
	return e, nil
}

// GetUser returns the User who submitted this job.
func (n *NewJob) GetUser() (*User, error) {
	u, err := user.Lookup(n.PBSJob.Origin.Logname)
	if err != nil {
		return nil, err
	}

	// TODO Go standard lib does not have a way to get primary groupname.
	// If Go starts supporting it this can be removed.
	cmd := exec.Command("id", "-g", "-n", u.Username)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	g := out.String()
	// Trim space and newline
	g = strings.TrimSpace(g)

	return &User{
		Username:  u.Username,
		UID:       u.Uid,
		Groupname: g,
		GID:       u.Gid,
		HomeDir:   u.HomeDir,
	}, nil
}

// GetShell returns the shell to use in the container.
func (n *NewJob) GetShell() (*Shell, error) {
	// If we are in interactive mode, don't tee.
	if n.PBSJob.Environment == "PBS_INTERACTIVE" {
		return &Shell{
			Stdout: os.Stdout,
			Stderr: os.Stderr,
			Stdin:  os.Stdin,
			TTY:    true,
		}, nil
	}

	return &Shell{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
		TTY:    false,
	}, nil
}

// GetImage returns a new Image.
func (n *NewJob) GetImage() (*Image, error) {

	if config.GetImageName() == "" && config.GetImageSource() == "" {
		return nil, fmt.Errorf("No image specifed.")
	}

	if config.GetImageName() != "" && config.GetImageSource() != "" {
		return nil, fmt.Errorf("Can not get image based on name and source.")
	}

	if config.GetImageName() != "" {
		return &Image{
			ImageName:   config.GetImageName(),
			ImageSource: config.GetImageSource(),
			Source:      false,
		}, nil
	}

	return &Image{
		ImageName:   config.GetImageName(),
		ImageSource: config.GetImageSource(),
		Source:      true,
		Privileged:  config.GetImagePrivileged(),
		InitCmd:     config.GetImageInit(),
	}, nil
}

// getMounts Get the default mounts
func (n *NewJob) addMounts(m Mounts) (Mounts, error) {
	// Not interactive job.
	if n.PBSJob.Environment != "PBS_INTERACTIVE" {
		s := fmt.Sprintf("%v%v.SC", config.GetJobDir(), n.PBSJob.JobID)
		m = m.AddMount(s, s, false) // Read only
	}

	// Mount home
	m = m.AddMount(n.PBSJob.Origin.Home, n.PBSJob.Origin.Home, true) // RW

	fmt.Println(m)

	return m, nil
}
