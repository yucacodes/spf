package main

import (
	"github.com/jessevdk/go-flags"
	"github.com/yucacodes/secure-port-forwarding/config"
	"github.com/yucacodes/secure-port-forwarding/node"
)

var NodeOptions struct {
	ConfigFile string `short:"c" long:"config" description:"Server config file" required:"true"`
}

func main() {
	_, err := flags.Parse(&NodeOptions)
	if err != nil {
		return
	}

	config, err := config.ConfigFromFile(NodeOptions.ConfigFile)
	if err != nil {
		panic("Invalid config")
	}
	node.NewNode(config).Run()
}
