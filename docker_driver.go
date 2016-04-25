package main

import (
	"io/ioutil"

	"fmt"

	"github.com/fsouza/go-dockerclient"
)

// DockerDriver driver for runnig job in Docker.
type DockerDriver struct {
	container *docker.Container
	client    *docker.Client
}

// NewDockerDriver returns a new docker driver.
func NewDockerDriver() (*DockerDriver, error) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return nil, err
	}

	return &DockerDriver{
		container: nil,
		client:    client,
	}, nil
}

func (d *DockerDriver) imageExists(i string) (bool, error) {

	images, err := d.client.ListImages(docker.ListImagesOptions{Filter: i})
	if err != nil {
		return false, err
	}

	for _, image := range images {
		for _, tag := range image.RepoTags {
			if tag == i {
				return true, nil
			}
		}
	}

	return false, nil
}

func (d *DockerDriver) startContainer() error {
	err := d.client.StartContainer(d.container.ID, d.container.HostConfig)
	if err != nil {
		return err
	}
	return nil
}

func (d *DockerDriver) removeContainer() error {
	err := d.client.RemoveContainer(docker.RemoveContainerOptions{
		ID:            d.container.ID,
		RemoveVolumes: true,
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *DockerDriver) waitContainer() error {
	exitCode, err := d.client.WaitContainer(d.container.ID)
	if err != nil {
		return err
	}

	if exitCode != 0 {
		err := fmt.Errorf("Prepare images exited with non-zero exit code: %v", exitCode)
		return err
	}

	return nil
}

func getPrepareContainerConfig(j *Job, c []string) *docker.Config {
	return &docker.Config{
		AttachStdin:  false,
		AttachStdout: false,
		AttachStderr: false,
		Tty:          false,
		OpenStdin:    false,
		StdinOnce:    false,
		Image:        j.Image.ImageName,
		Cmd:          c,
	}
}

// Prepare Implementation of the prepare interface.
func (d *DockerDriver) Prepare(j *Job) error {

	// Check if image exist
	hasImage, err := d.imageExists(j.Image.ImageName)
	if err != nil {
		return err
	}

	// Pull base image
	if !hasImage {
		err = d.client.PullImage(docker.PullImageOptions{
			Repository:   j.Image.ImageName,
			OutputStream: ioutil.Discard,
		}, docker.AuthConfiguration{})
		if err != nil {
			return err
		}
	}

	groupAdd := fmt.Sprintf("groupadd -g %v %v", j.User.GID, j.User.Groupname)
	addUser := fmt.Sprintf("adduser -d %v -u %v -g %v %v", j.User.HomeDir, j.User.UID, j.User.GID, j.User.Username)

	cmd := []string{
		"/bin/bash",
		"-c",
		groupAdd + " && " + addUser,
	}

	d.container, err = d.client.CreateContainer(docker.CreateContainerOptions{
		Config:     getPrepareContainerConfig(j, cmd),
		HostConfig: &docker.HostConfig{},
	})
	if err != nil {
		return err
	}

	err = d.startContainer()
	if err != nil {
		return err
	}

	err = d.waitContainer()
	if err != nil {
		return err
	}

	image, err := d.client.CommitContainer(docker.CommitContainerOptions{Container: d.container.ID})
	if err != nil {
		return err
	}

	err = d.removeContainer()
	if err != nil {
		return err
	}

	// Prepare the run container.
	d.container, err = d.client.CreateContainer(docker.CreateContainerOptions{
		Config:     getRunContainerConfig(j, image),
		HostConfig: getRunHostConfig(j),
	})
	if err != nil {
		return err
	}

	return nil
}

// Run implementation of the driver interface.
func (d *DockerDriver) Run(j *Job) error {

	err := d.startContainer()
	if err != nil {
		return err
	}

	err = d.client.AttachToContainer(docker.AttachToContainerOptions{
		Container:    d.container.ID,
		InputStream:  j.Shell.Stdin,
		OutputStream: j.Shell.Stdout,
		ErrorStream:  j.Shell.Stderr,
		Logs:         true,
		Stream:       true,
		Stdin:        true,
		Stdout:       true,
		Stderr:       true,
		RawTerminal:  j.Shell.TTY, // Use raw terminal with tty https://godoc.org/github.com/fsouza/go-dockerclient#AttachToContainerOptions
	})
	if err != nil {
		return err
	}

	err = d.removeContainer()
	if err != nil {
		return err
	}

	return nil
}

func getRunContainerConfig(j *Job, i *docker.Image) *docker.Config {
	return &docker.Config{
		User:         j.User.Username,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          j.Shell.TTY,
		OpenStdin:    true,
		StdinOnce:    true,
		Env:          j.Env,
		Image:        i.ID,
		Cmd:          j.Cmd,
	}
}

func getRunHostConfig(j *Job) *docker.HostConfig {
	var binds []string
	for _, m := range j.Mounts {
		binds = append(binds, m.DockerBindString())
	}
	return &docker.HostConfig{Binds: binds}
}
