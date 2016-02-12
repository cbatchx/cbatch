package main

import (
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
	Container      *docker.Container `json:"-"`
	User           *user.User        `json:"-"`
}

// GetScript gets the Script of the job
func (j *Job) GetScript() *Script {
	return &Script{jobFolder + j.ID + scriptFileEnding, nil}
}

func (j *Job) setUser() error {
	u, err := user.Lookup(j.Username)
	if err != nil {
		return err
	}

	j.User = u
	return nil
}

// CreateRunImage creates a image for basic running of the image
func (j *Job) CreateRunImage() error {
	err := j.setUser()
	if err != nil {
		return err
	}

	client, err := docker.NewClientFromEnv()
	if err != nil {
		return err
	}

	// Dont set this on the user only used to create new image
	c, err := client.CreateContainer(docker.CreateContainerOptions{
		Config:     j.getUserInsertContainerConfig(),
		HostConfig: getDefaultHostConfig(),
	})
	if err != nil {
		return err
	}
	log.Printf("Created temporary container %v \n", c)

	err = client.StartContainer(c.ID, c.HostConfig)
	if err != nil {
		return err
	}
	log.Printf("Started temporary container %v \n", c)

	// img, err := client.CommitContainer(docker.CommitContainerOptions{
	//	Container: c.ID,
	//	Run:       j.getDefaultContainerConfig(),
	//})
	//	if err != nil {
	//		return err
	//	}

	//log.Printf("Created new image %v \n", img)

	return nil
}

// StartContainer Start a container for this job
func (j *Job) StartContainer() error {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return err
	}

	c, err := client.CreateContainer(docker.CreateContainerOptions{
		Name:       "",
		Config:     j.getDefaultContainerConfig(),
		HostConfig: getDefaultHostConfig(),
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

// DecodeJob Decodes a job
func DecodeJob(r io.Reader) (*Job, error) {
	var j Job
	dec := json.NewDecoder(r)
	err := dec.Decode(&j)
	return &j, err
}

func (j *Job) getUserInsertContainerConfig() *docker.Config {
	groupadd := fmt.Sprintf("groupadd -f -g %v %v", j.User.Gid, j.Group)
	useradd := fmt.Sprintf("useradd -u %v -g %v %v", j.User.Uid, j.Group, j.User.Username)
	mkdir := fmt.Sprintf("mkdir --parent %v", j.User.HomeDir)
	chown := fmt.Sprintf("chown -R %v:%v %v", j.User.Username, j.Group, j.User.HomeDir)
	cmd := fmt.Sprintf("/bin/bash -c \" %v && %v && %v && %v", groupadd, useradd, mkdir, chown)

	return &docker.Config{
		AttachStdin:  false,
		AttachStdout: false,
		AttachStderr: false,
		Tty:          false,
		OpenStdin:    false,
		StdinOnce:    false,
		Image:        "centos:latest",
		Cmd:          []string{cmd},
	}
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
		Image:        "centos:latest",
		Cmd:          []string{"tail -f /dev/null"},
	}
}

func getDefaultHostConfig() *docker.HostConfig {
	return &docker.HostConfig{}
}
