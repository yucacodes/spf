package config

import "strconv"

type Node struct {
	Name    string   `yaml:"name" json:"name"`
	Host    *string  `yaml:"host" json:"host"`
	Port    *int     `yaml:"port" json:"port"`
	Key     *string  `yaml:"key" json:"key"`
	Through *Through `yaml:"through" json:"through"`
}

const DefaultNodePort = 5000

func (n *Node) Connection() string {
	if n.Host == nil {
		return ""
	}
	host := *n.Host
	port := strconv.Itoa(DefaultNodePort)
	if n.Port != nil {
		port = strconv.Itoa(*n.Port)
	}
	return host + ":" + port
}

func (n *Node) IsPublic() bool {
	return n.Host != nil
}
