package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func LoadConfig(configPath string) (*Properties, error) {
	if configPath == "" {
		return nil, errors.New("config is required")
	}
	return loadConfig(configPath)
}

func loadConfig(configPath string) (*Properties, error) {
	configDir, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(configDir)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	// expand env vars for secrets
	b = []byte(os.ExpandEnv(string(b)))

	cfg, err := parseConfig(b)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func parseConfig(b []byte) (*Properties, error) {
	cfg := new(Properties)
	if err := yaml.Unmarshal(b, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
