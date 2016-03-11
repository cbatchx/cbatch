package main

import "github.com/BurntSushi/toml"

// Config represents the config.
type Config struct {
	MOMPriv string
}

// ReadConfig reads the config and returns it.
func ReadConfig() (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
