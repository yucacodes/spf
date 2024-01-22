package client

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type ClientConfig struct {
	ServerHost string `yaml:"serverHost" json:"serverHost"`
	ServerPort int    `yaml:"serverPort" json:"serverPort"`
	AppKey     string `yaml:"appKey" json:"appKey"`
	AppHost    string `yaml:"appHost" json:"appHost"`
	AppPort    int    `yaml:"appPort" json:"appPort"`
}

func ClientConfigFromFile(path string) (*ClientConfig, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := ClientConfig{}

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
