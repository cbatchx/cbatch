package main

import "github.com/BurntSushi/toml"

// Config represents the config.
type Config struct {
	MOMPriv string `toml:"mom_priv"`
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
	return c.MOMPriv
}
