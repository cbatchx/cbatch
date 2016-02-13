package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os/user"

	"github.com/dizk/docquer/Godeps/_workspace/src/github.com/fsouza/go-dockerclient"
)

const (
	jobFolder        = mompriv + "jobs/"
	scriptFileEnding = ".SC"
)

// Job represents a job
// Example JSON of job
//{
//	"jobid": "36.master",
//	"user": "vagrant",
//	"group": "vagrant",
//	"jobname": "STDIN",
//	"resourcelimits": "neednodes=1,nodes=1,walltime=01:00:00",
//	"jobqueue": "batch",
//	"jobaccount": ""
// }
type Job struct {
	ID             string            `json:"jobid"`
	Username       string            `json:"user"`
	Group          string            `json:"group"`
	Name           string            `json:"jobname"`
	ResourceLimits string            `json:"resourcelimits"`
	Queue          string            `json:"jobqueue"`
	Account        string            `json:"jobaccount"`
	Cmd            string            `json:"cmd"`
	ImageName      string            `json:"-"`
	ScriptPath     string            `json:"-"`
	Container      *docker.Container `json:"-"`
	User           *user.User        `json:"-"`
}

// InitJob Initialize job
func (j *Job) InitJob() error {
	// Set the path to the script
	j.ScriptPath = jobFolder + j.ID + scriptFileEnding
	// Set User struct
	err := j.setUser()
	if err != nil {
		return err
	}

	return nil
}

func (j *Job) setUser() error {
	u, err := user.Lookup(j.Username)
	if err != nil {
		return err
	}

	j.User = u
	return nil
}

// CreateImage creates a image for basic running of the image
func (j *Job) CreateImage() error {

	groupadd := fmt.Sprintf("groupadd -f -g %v %v", j.User.Gid, j.Group)
	useradd := fmt.Sprintf("useradd -u %v -g %v %v", j.User.Uid, j.Group, j.User.Username)
	mkdir := fmt.Sprintf("mkdir --parent %v", j.User.HomeDir)
	chown := fmt.Sprintf("chown -R %v:%v %v", j.User.Username, j.Group, j.User.HomeDir)

	d := DockerFile{
		From: j.ImageName,
		Run:  []string{groupadd, useradd, mkdir, chown},
		User: j.User.Username,
		Cmd:  "/bin/bash -c tail -f /dev/null",
	}

	in, err := GetTarBuf(d)
	if err != nil {
		return err
	}
	out := new(bytes.Buffer)

	client, err := docker.NewClientFromEnv()
	if err != nil {
		return err
	}

	err = client.BuildImage(docker.BuildImageOptions{
		Name:         j.ID,
		InputStream:  in,
		OutputStream: out,
		Pull:         true,
	})
	if err != nil {
		return err
	}

	j.ImageName = j.ID

	return nil
}

// StartContainer Start a container for this job
func (j *Job) StartContainer() error {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return err
	}

	c, err := client.CreateContainer(docker.CreateContainerOptions{
		Config:     j.getDefaultContainerConfig(),
		HostConfig: j.getDefaultHostConfig(),
	})
	if err != nil {
		return err
	}

	j.Container = c
	err = client.StartContainer(j.Container.ID, j.Container.HostConfig)
	if err != nil {
		return err
	}

	return nil
}

// StopContainer a
func (j *Job) StopContainer() error {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return err
	}

	err = client.StopContainer(j.Container.ID, 0)
	if err != nil {
		return err
	}

	return nil
}

// RemoveContainer deletes a container
func (j *Job) RemoveContainer() error {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return err
	}

	err = client.RemoveContainer(docker.RemoveContainerOptions{
		ID:            j.Container.ID,
		RemoveVolumes: true,
	})
	if err != nil {
		return err
	}

	return nil
}

// RemoveImage deletes an image
func (j *Job) RemoveImage() error {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return err
	}

	err = client.RemoveImage(j.ImageName)
	if err != nil {
		return err
	}

	return nil
}

// DecodeJob Decodes a job
func DecodeJob(r io.Reader) (Job, error) {
	var j Job
	dec := json.NewDecoder(r)
	err := dec.Decode(&j)
	return j, err
}

func (j *Job) getDefaultContainerConfig() *docker.Config {
	return &docker.Config{
		User:         j.User.Username,
		AttachStdin:  false,
		AttachStdout: false,
		AttachStderr: false,
		Tty:          false,
		OpenStdin:    false,
		StdinOnce:    false,
		Image:        j.ImageName,
		Cmd:          []string{"/bin/bash", "-c", "tail -f /dev/null"},
		// Mounts:       []docker.Mount{j.getScriptMount()},
	}
}

func (j *Job) getDefaultHostConfig() *docker.HostConfig {
	scriptBind := fmt.Sprintf("%v:%v", j.ScriptPath, j.ScriptPath)
	hc := docker.HostConfig{Binds: []string{scriptBind}}
	log.Printf("[JOB: %v] Using hostconfig: %v\n", j.ID, hc)
	return &hc
}

func (j *Job) getScriptMount() docker.Mount {
	return docker.Mount{
		Source:      j.ScriptPath,
		Destination: j.ScriptPath,
	}
}
