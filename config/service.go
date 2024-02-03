package config

import "strconv"

type Service struct {
	Name    string   `yaml:"name" json:"name"`
	Host    *string  `yaml:"host" json:"host"`
	Port    *int     `yaml:"port" json:"port"`
	Through *Through `yaml:"through" json:"through"`
}

const DefaultServiceHost = "127.0.0.1"

func (n *Service) Connection() string {
	if n.Port == nil {
		return ""
	}
	port := strconv.Itoa(*n.Port)
	host := DefaultServiceHost
	if n.Host != nil {
		host = *n.Host
	}
	return host + ":" + port
}

func (n *Service) IsOwn() bool {
	return n.Connection() != ""
}
