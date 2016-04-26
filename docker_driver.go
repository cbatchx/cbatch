package main

import (
	"io/ioutil"

	"fmt"

	"net/url"

	"log"

	"os"

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

// Prepare Implementation of the prepare interface.
func (d *DockerDriver) Prepare(j *Job) error {

	imageName, err := d.getImage(j)
	if err != nil {
		return err
	}

	groupAdd := fmt.Sprintf("groupadd -g %v %v", j.User.GID, j.User.Groupname)
	addUser := fmt.Sprintf("adduser -d %v -u %v -g %v %v", j.User.HomeDir, j.User.UID, j.User.GID, j.User.Username)

	content := []byte("#!/bin/bash\n" + groupAdd + "\n" + addUser)
	tmpfile, err := ioutil.TempFile("", "cbatch_prepare")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

	cmd := []string{
		"/bin/bash",
		"/prepare.sh",
	}

	// Prepend the init command
	if j.Image.InitCmd != "" {
		cmd = append([]string{j.Image.InitCmd}, cmd...)
	}

	log.Println(cmd)

	m := Mount{
		Source:      tmpfile.Name(),
		Destination: "/prepare.sh",
		RW:          true,
	}

	d.container, err = d.client.CreateContainer(docker.CreateContainerOptions{
		Config:     getPrepareContainerConfig(imageName, cmd),
		HostConfig: getPrepareHostConfig(j, m),
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
		HostConfig: getHostConfig(j),
	})

	return err
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

	return err
}

func (d *DockerDriver) getImage(j *Job) (string, error) {
	if j.Image.Source {
		imageName, err := d.importImage(j.Image.ImageSource)
		return imageName, err
	}
	err := d.pullImage(j.Image.ImageName)
	return j.Image.ImageName, err

}

func (d *DockerDriver) importImage(imageSource string) (string, error) {

	// Use path as image name
	u, err := url.Parse(imageSource)
	if err != nil {
		return "", err
	}
	imageName := u.Path[1:]

	hasImage, err := d.imageExists(imageName + ":latest")
	if err != nil {
		return "", err
	}

	if hasImage {
		return imageName, nil
	}

	err = d.client.ImportImage(docker.ImportImageOptions{
		Source:       imageSource,
		Repository:   imageName,
		OutputStream: ioutil.Discard,
	})
	return imageName, err
}

func (d *DockerDriver) pullImage(imageName string) error {
	// Check if image exist
	hasImage, err := d.imageExists(imageName)
	if err != nil {
		return err
	}

	if hasImage {
		return nil
	}

	// Pull image
	err = d.client.PullImage(docker.PullImageOptions{
		Repository:   imageName,
		OutputStream: ioutil.Discard,
	}, docker.AuthConfiguration{})

	return err
}

func (d *DockerDriver) startContainer() error {
	err := d.client.StartContainer(d.container.ID, d.container.HostConfig)
	return err
}

func (d *DockerDriver) waitContainer() error {
	exitCode, err := d.client.WaitContainer(d.container.ID)
	if err != nil {
		return err
	}

	if exitCode != 0 {
		err := fmt.Errorf("Image preparation exited with non-zero exit code: %v", exitCode)
		return err
	}

	return nil
}

func (d *DockerDriver) removeContainer() error {
	err := d.client.RemoveContainer(docker.RemoveContainerOptions{
		ID:            d.container.ID,
		RemoveVolumes: true,
	})
	return err
}

func getPrepareContainerConfig(imageName string, cmd []string) *docker.Config {
	return &docker.Config{
		AttachStdin:  false,
		AttachStdout: false,
		AttachStderr: false,
		Tty:          false,
		OpenStdin:    false,
		StdinOnce:    false,
		Image:        imageName,
		Cmd:          cmd,
	}
}

func getPrepareHostConfig(j *Job, m Mount) *docker.HostConfig {
	var binds []string
	binds = append(binds, m.DockerBindString())

	return &docker.HostConfig{
		Binds:      binds,
		Privileged: j.Image.Privileged,
	}
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

func getHostConfig(j *Job) *docker.HostConfig {
	var binds []string
	for _, m := range j.Mounts {
		binds = append(binds, m.DockerBindString())
	}
	return &docker.HostConfig{
		Binds:      binds,
		Privileged: j.Image.Privileged,
	}
}

func (d *DockerDriver) imageExists(i string) (bool, error) {

	images, err := d.client.ListImages(docker.ListImagesOptions{})
	if err != nil {
		return false, err
	}

	for _, image := range images {
		for _, tag := range image.RepoTags {
			if tag == i {
				log.Println(tag)
				return true, nil
			}
		}
	}

	return false, nil
}
