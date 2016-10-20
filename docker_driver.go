package main

import (
	"io/ioutil"
	"time"

	"net/url"

	log "github.com/Sirupsen/logrus"

	"github.com/fsouza/go-dockerclient"
)

// DockerDriver driver for runnig job in Docker.
type DockerDriver struct {
	container *docker.Container
	client    *docker.Client
}

// NewDockerDriver returns a new docker driver.
func NewDockerDriver() (*DockerDriver, error) {
	client, err := docker.NewVersionedClientFromEnv("") // TODO define version in config
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
	defer MeasureTime(time.Now(), log.Fields{"job": j, "job_id": j.ID}, "Container prepared.")

	imageName, err := d.getImage(j)
	if err != nil {
		return err
	}

	// Prepare the run container.
	d.container, err = d.client.CreateContainer(docker.CreateContainerOptions{
		Config:     getRunContainerConfig(j, imageName),
		HostConfig: getHostConfig(j),
	})

	return err
}

// Run implementation of the driver interface.
func (d *DockerDriver) Run(j *Job) error {
	defer MeasureTime(time.Now(), log.Fields{"job": j, "job_id": j.ID}, "Job finished.")

	err := d.startContainer()
	if err != nil {
		log.Warn("Failed to start the container.")
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
		log.Fatal("Failed to attach to the container.")
		return err
	}

	err = d.removeContainer()

	if err != nil {
		log.Fatal("Failed to remove the container.")
	}

	return err
}

func (d *DockerDriver) Abort() error {
	// No container created nothing to do
	if d.container == nil {
		return nil
	}
	// Remove Container and Force kill if needed.
	return d.client.RemoveContainer(docker.RemoveContainerOptions{
		ID:            d.container.ID,
		Force:         true,
		RemoveVolumes: true,
	})
}

func (d *DockerDriver) startContainer() error {
	err := d.client.StartContainer(d.container.ID, d.container.HostConfig)
	return err
}

func (d *DockerDriver) removeContainer() error {
	err := d.client.RemoveContainer(docker.RemoveContainerOptions{
		ID:            d.container.ID,
		RemoveVolumes: true,
	})
	return err
}

func getRunContainerConfig(j *Job, imageName string) *docker.Config {
	return &docker.Config{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          j.Shell.TTY,
		OpenStdin:    true,
		StdinOnce:    true,
		Env:          j.Env,
		Image:        imageName,
		Cmd:          j.Cmd,
	}
}

func getHostConfig(j *Job) *docker.HostConfig {
	return &docker.HostConfig{
		Binds:      buildBindString(j.Mounts),
		Privileged: j.Image.Privileged,
	}
}

func buildBindString(mounts Mounts) []string {
	var binds []string
	for _, mount := range mounts {
		binds = append(binds, mount.DockerBindString())
	}
	return binds
}

func (d *DockerDriver) getImage(j *Job) (string, error) {
	if j.Image.Source {
		log.WithFields(log.Fields{"image": j.Image}).Info("Dowloading image from source.")
		imageName, err := d.importImage(j.Image.ImageSource)
		return imageName, err
	}

	log.WithFields(log.Fields{"image": j.Image}).Info("Pulling image from docker hub.")
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
