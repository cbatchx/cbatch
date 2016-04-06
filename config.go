package main

import "github.com/BurntSushi/toml"

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
	Name   string
	Source string
}

type influxConfig struct {
	Host     string
	User     string
	Password string
	Present  bool
}

const configPath = "/etc/cbatch.toml"

// ReadConfig reads the config and returns it.
func ReadConfig(c *Config) error {
	md, err := toml.DecodeFile(configPath, c)
	if err != nil {
		return err
	}

	if md.IsDefined("influx") {
		c.Influx.Present = true
	} else {
		c.Influx.Present = false
	}

	return nil
}

// GetMOMPriv get mom_priv folder from config.
func (c *Config) GetMOMPriv() string {
	return c.Torque.MOMPriv
}

// GetImageName get the name of the image to run.
func (c *Config) GetImageName() string {
	return c.Image.Name
}

// GetImageSource get the server to download the image from.
func (c *Config) GetImageSource() string {
	return ""
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
