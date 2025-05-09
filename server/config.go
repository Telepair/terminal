// Package server provides HTTP server and configuration utilities.
package server

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds server configuration.
type Config struct {
	Addr string `yaml:"addr"`
}

// LoadConfig loads configuration from the given YAML file path.
func LoadConfig(path string) (*Config, error) {
	// #nosec G304 -- Configuration file path is controlled by trusted user input.
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close() // Ignore close error
	}()
	var cfg Config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
