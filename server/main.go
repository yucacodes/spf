package server

import (
	"fmt"
	"net"
	"strconv"

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
	}
	runServer(*config)
}

func runServer(config ServerConfig) {
	appsListener, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(config.Port))

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer appsListener.Close()

	fmt.Println("Server is listening on port " + strconv.Itoa(config.Port))
	listenApps(config.Apps, appsListener)
}
