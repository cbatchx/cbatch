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

	return &Job{
		User:   u,
		Cmd:    n.Args,
		Shell:  s,
		Image:  i,
		Mounts: m,
	}, nil
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

// GetShell returns the shell and optinally an error
// It also returns an stdin io.Reader, enabling to read the commands piped to the final shell.
// It creates a tee reader so the Stdin is still passed to the driver.
// If the job is interactive the tee reader will not be set.
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
	return &Image{
		ImageName:   "centos:latest",
		ImageSource: "",
	}, nil
}

// getMounts Get the default mounts
func (n *NewJob) addMounts(m Mounts) (Mounts, error) {
	// Not interactive job.
	if n.PBSJob.Environment != "PBS_INTERACTIVE" {
		s := fmt.Sprintf("%v%v.SC", config.GetMOMPriv(), n.PBSJob.JobID)
		m = m.AddMount(s, s, false) // Read only
	}

	// Mount home
	m = m.AddMount(n.PBSJob.Origin.Home, n.PBSJob.Origin.Home, true) // RW

	// Mount /etc/passwd and /etc/group
	// To allow for other users to read
	m = m.AddMount("/etc/passwd", "/etc/passwd", false) // RO
	m = m.AddMount("/etc/group", "/etc/group", false)   // RO

	fmt.Println(m)

	return m, nil
}
