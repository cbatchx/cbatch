package main

import "github.com/fsouza/go-dockerclient"

// DockerDriver driver for runnig job in Docker.
type DockerDriver struct {
	container *docker.Container
}

// NewDockerDriver returns a new docker driver.
func NewDockerDriver() *DockerDriver {
	return &DockerDriver{container: nil}
}

// Run implementation of the driver interface.
func (d *DockerDriver) Run(j *Job) error {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return err
	}

	// Pull image
	// err = client.PullImage(docker.PullImageOptions{
	// 	Repository:   j.Image.ImageName,
	// 	OutputStream: ioutil.Discard,
	// }, docker.AuthConfiguration{})
	// if err != nil {
	// 	return err
	// }

	d.container, err = client.CreateContainer(docker.CreateContainerOptions{
		Config:     getDefaultContainerConfig(j),
		HostConfig: getDefaultHostConfig(j),
	})
	if err != nil {
		return err
	}

	err = client.StartContainer(d.container.ID, d.container.HostConfig)
	if err != nil {
		return err
	}

	err = client.AttachToContainer(docker.AttachToContainerOptions{
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

	// Since StdinOnce is set the container will stop automatically

	err = client.RemoveContainer(docker.RemoveContainerOptions{
		ID:            d.container.ID,
		RemoveVolumes: true,
	})
	if err != nil {
		return err
	}

	return nil
}

func getDefaultContainerConfig(j *Job) *docker.Config {
	return &docker.Config{
		User:         j.User.Username,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          j.Shell.TTY,
		OpenStdin:    true,
		StdinOnce:    true,
		Image:        j.Image.ImageName,
		Cmd:          j.Cmd,
	}
}

func getDefaultHostConfig(j *Job) *docker.HostConfig {
	var binds []string
	for _, m := range j.Mounts {
		binds = append(binds, m.DockerBindString())
	}
	return &docker.HostConfig{Binds: binds}
}
