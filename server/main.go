package server

import (
	"log"
	"os"

	"github.com/jessevdk/go-flags"
)

var ServerOptions struct {
	ConfigFile string `short:"c" long:"config" description:"Server config file" required:"true"`
}

func Main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	_, err := flags.Parse(&ServerOptions)
	if err != nil {
		return
	}

	config, err := ServerConfigFromFile(ServerOptions.ConfigFile)
	if err != nil {
		logger.Fatalln(err)
		return
	}
	server := NewServer(config.Port, config.Apps)
	server.Listen()
}
