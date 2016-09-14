package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"

	log "github.com/Sirupsen/logrus"
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

	c, err := n.GetCmd(i)
	if err != nil {
		return nil, err
	}

	job := &Job{
		User:   u,
		Cmd:    c,
		Shell:  s,
		Image:  i,
		Mounts: m,
		Env:    e,
	}

	bs, err := NewBootstrap(job)
	if err != nil {
		return nil, err
	}

	job.Mounts = job.Mounts.AddMount(bs.GetScriptPath(), "/bootstrap.sh", false)

	log.WithFields(log.Fields{
		"job": job,
	}).Info("Parsed job")

	return job, nil
}

// GetCmd get the command to run in the container
func (n *NewJob) GetCmd(i *Image) ([]string, error) {
	cmd := n.Args

	// Prepend the bootstrap script
	cmd = append([]string{"/bootstrap.sh"}, cmd...)

	// Prepend the init command
	if i.InitCmd != "" {
		cmd = append([]string{i.InitCmd}, cmd...)
	}

	return cmd, nil
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

	source := false

	if config.GetImageSource() != "" {
		source = true
	}

	return &Image{
		ImageName:   config.GetImageName(),
		ImageSource: config.GetImageSource(),
		Source:      source,
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
	// m = m.AddMount(n.PBSJob.Origin.Home, n.PBSJob.Origin.Home, true) // RW

	// Mount cvmfs if present
	if config.GetCvmfs() != "" {
		m = m.AddMount(config.GetCvmfs(), "/cvmfs", true)
	}

	return m, nil
}
