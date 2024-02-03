package config

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

type NodeId struct {
	Name string `yaml:"name" json:"name"`
	Key  string `yaml:"key" json:"key"`
}

type Through struct {
	Node string `yaml:"node" json:"node"`
}

type Publish struct {
	Service string   `yaml:"service" json:"service"`
	Through *Through `yaml:"through" json:"through"`
}

type Connect struct {
	Service string `yaml:"service" json:"service"`
}

type Listen struct {
	Port    int     `yaml:"port" json:"port"`
	Connect Connect `yaml:"connect" json:"connect"`
}

func (l *Listen) ListenConnection() string {
	host := "0.0.0.0"
	port := strconv.Itoa(l.Port)
	return host + ":" + port
}

type Config struct {
	Port                  *int       `yaml:"port" json:"port"`
	Id                    *NodeId    `yaml:"id" json:"id"`
	Nodes                 []*Node    `yaml:"nodes" json:"nodes"`
	Services              []*Service `yaml:"services" json:"services"`
	Publish               []*Publish `yaml:"publish" json:"publish"`
	Listen                []*Listen  `yaml:"listen" json:"listen"`
	DisableNodeValidation *bool      `yaml:"disableNodeValidation" json:"disableNodeValidation"`
}

func (c *Config) ListenConnection() string {
	host := "0.0.0.0"
	port := strconv.Itoa(DefaultNodePort)
	if c.Port != nil {
		port = strconv.Itoa(*c.Port)
	}
	return host + ":" + port
}

func ConfigFromFile(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := Config{}

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
