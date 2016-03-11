package main

import "github.com/BurntSushi/toml"

// Config represents the config.
type Config struct {
	Torque torqueConfig
	Image  imageConfig
}

type torqueConfig struct {
	MOMPriv string `toml:"mom_priv"`
}

type imageConfig struct {
	Name   string
	Source string
}

const configPath = "/etc/cbatch.toml"

// ReadConfig reads the config and returns it.
func ReadConfig(c *Config) error {
	if _, err := toml.DecodeFile(configPath, c); err != nil {
		return err
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
