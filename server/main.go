package server

import (
	"fmt"

	"github.com/jessevdk/go-flags"
)

var ServerOptions struct {
	ConfigFile string `short:"c" long:"config" description:"Server config file" required:"true"`
}

func Main() {
	_, err := flags.Parse(&ServerOptions)
	if err != nil {
		return
	}

	config, err := ServerConfigFromFile(ServerOptions.ConfigFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	server := NewServer(config.Port, config.Apps)
	defer server.Close()
	server.Listen()
}
