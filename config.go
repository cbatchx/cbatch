package main

import (
	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
)

// Config represents the config.
type Config struct {
	Torque torqueConfig
	Image  imageConfig
	Influx influxConfig
}

type torqueConfig struct {
	MOMPriv string `toml:"mom_priv"`
}

type imageConfig struct {
	Name       string
	Source     string
	Privileged bool
	Init       string
	Cvmfs      string
	MountHome  bool `toml:"mount_home"`
	MountHosts bool `toml:"mount_hosts"`
}

type influxConfig struct {
	Host     string
	User     string
	Password string
	Present  bool
}

const configPath = "/var/lib/cbatch/config/config.toml"

const configFallback = "config/config.toml"

// ReadConfig reads the config and returns it.
func ReadConfig(c *Config) error {

	md, err := toml.DecodeFile(configPath, c)
	if err != nil {
		log.Error(err)

		log.Warn("Falling back to: ", configFallback)
		md, err = toml.DecodeFile(configFallback, c)

		if err != nil {
			log.Fatal("Failed to read fallback config. ", err)
		}

	}

	if md.IsDefined("influx") {
		c.Influx.Present = true
	} else {
		c.Influx.Present = false
	}

	log.WithFields(log.Fields{
		"config": c,
	}).Info("Read config")
	return nil
}

// GetMOMPriv get mom_priv folder from config.
func (c *Config) GetMOMPriv() string {
	return c.Torque.MOMPriv
}

// GetJobDir get the folder where jobs are located. Usually mom_priv/jobs
func (c *Config) GetJobDir() string {
	return c.Torque.MOMPriv + "jobs/"
}

// GetImageName get the name of the image to run.
func (c *Config) GetImageName() string {
	return c.Image.Name
}

// GetImageSource get the server to download the image from.
func (c *Config) GetImageSource() string {
	return c.Image.Source
}

// GetImagePrivileged whetever to run image as Privileged or not.
func (c *Config) GetImagePrivileged() bool {
	return c.Image.Privileged
}

// GetImageInit get special command to run before running the job.
func (c *Config) GetImageInit() string {
	return c.Image.Init
}

// GetCvmfs get the path of Cvmfs
func (c *Config) GetCvmfs() string {
	return c.Image.Cvmfs
}

// MountHome returns a boolean to mount the home folder or not
func (c *Config) MountHome() bool {
	return c.Image.MountHome
}

// MountHosts returns a boolean to mount /etc/hosts or not
func (c *Config) MountHosts() bool {
	return c.Image.MountHosts
}

// InfluxAvailable check if influxdb is configured.
func (c *Config) InfluxAvailable() bool {
	return c.Influx.Present
}

// GetInfluxHost get the host of the influx database
// Returns "" if influx is not configured
func (c *Config) GetInfluxHost() string {
	if !c.Influx.Present {
		return ""
	}
	return c.Influx.Host
}

// GetInfluxUser get the user of the influx database
func (c *Config) GetInfluxUser() string {
	if !c.Influx.Present {
		return ""
	}
	return c.Influx.User
}

// GetInfluxPassword get the host of the influx database
func (c *Config) GetInfluxPassword() string {
	if !c.Influx.Present {
		return ""
	}
	return c.Influx.Password
}

// GetConfigPath return path where config was read from.
func (c *Config) GetConfigPath() string {
	return configPath
}
