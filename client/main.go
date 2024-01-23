package client

import (
	"log"
	"os"

	"github.com/jessevdk/go-flags"
)

var ClientOptions struct {
	ConfigFile string `short:"c" long:"config" description:"Server config file" required:"true"`
}

func Main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	_, err := flags.Parse(&ClientOptions)
	if err != nil {
		return
	}

	config, err := ClientConfigFromFile(ClientOptions.ConfigFile)
	if err != nil {
		logger.Fatalln(err)
		return
	}
	client := NewClient(config.ServerHost, config.ServerPort, config.AppKey, config.AppHost, config.AppPort)
	client.Connect()
}
