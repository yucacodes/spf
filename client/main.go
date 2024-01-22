package client

import (
	"fmt"

	"github.com/jessevdk/go-flags"
)

var ClientOptions struct {
	ConfigFile string `short:"c" long:"config" description:"Server config file" required:"true"`
}

func Main() {
	_, err := flags.Parse(&ClientOptions)
	if err != nil {
		return
	}

	config, err := ClientConfigFromFile(ClientOptions.ConfigFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	client := NewClient(config.ServerHost, config.ServerPort, config.AppKey, config.AppHost, config.AppPort)
	defer client.Close()
	client.Connect()
}
