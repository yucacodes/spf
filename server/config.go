package server

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Port int          `yaml:"port" json:"port"`
	Apps []*AppConfig `yaml:"apps" json:"apps"`
}

func ServerConfigFromFile(path string) (*ServerConfig, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := ServerConfig{}

	if strings.HasSuffix(path, "yml") || strings.HasSuffix(path, "yaml") {
		err = yaml.Unmarshal(file, &config)
	} else if strings.HasSuffix(path, "json") {
		err = json.Unmarshal(file, &config)
	} else {
		return nil, errors.New("config file format not supported")
	}
	if err != nil {
		return nil, err
	}

	return &config, nil
}
