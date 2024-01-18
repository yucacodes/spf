package server

import (
	"net"
	"strconv"
)

type AppConfig struct {
	Port int    `yaml:"port" json:"port"`
	Key  string `yaml:"key" json:"key"`
}

type Server struct {
	port        int
	appsConfigs map[string]*AppConfig
}

func NewServer(port int) *Server {
	o := Server{port: port}
	return &o
}

func (s *Server) Listen() error {
	server, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(s.port))
	if err != nil {
		return err
	}
	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			return nil
		}

		// TODO

	}

}
