package server

import (
	"errors"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/yucacodes/secure-port-forwarding/app"
	"github.com/yucacodes/secure-port-forwarding/socket"
)

type AppConfig struct {
	Port int    `yaml:"port" json:"port"`
	Key  string `yaml:"key" json:"key"`
}

type Server struct {
	port        int
	appsConfigs map[string]*AppConfig
	appsServers map[string]*app.AppServer
	logger      *log.Logger
}

func NewServer(port int, appsConfigs []*AppConfig) *Server {
	s := Server{
		port:        port,
		appsConfigs: make(map[string]*AppConfig),
		logger:      log.New(os.Stdout, "Server: ", log.Ldate|log.Ltime),
	}

	for _, appConfig := range appsConfigs {
		s.appsConfigs[appConfig.Key] = appConfig
	}
	return &s
}

func (s *Server) Listen() error {
	server, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(s.port))
	if err != nil {
		s.logger.Fatal(err)
		return err
	}
	defer server.Close()
	s.logger.Println("Listening on port " + strconv.Itoa(s.port))
	for {
		conn, err := server.Accept()
		if err != nil {
			s.logger.Fatalln(err)
			return nil
		}
		s.logger.Println("New connection")
		jSocket := socket.NewJsonSocket(conn)
		req, err := s.GetAppRequest(jSocket)
		if err != nil {
			jSocket.Close()
			continue
		}

		if req.InitApp {
			go s.StartApp(jSocket, req.AppKey)
		} else if req.BackendToAppClient {
			s.RegisterAppClientBackend(jSocket, req.AppKey, req.AppClientId)
		}
	}

}

type AppRequest struct {
	AppKey             string `yaml:"appKey" json:"appKey"`
	InitApp            bool   `yaml:"initApp" json:"initApp"`
	BackendToAppClient bool   `yaml:"backendToAppClient" json:"backendToAppClient"`
	AppClientId        string `yaml:"appClientId" json:"appClientId"`
}

func (s *Server) GetAppRequest(jSocket *socket.JsonSocket) (*AppRequest, error) {
	req := AppRequest{}
	err := jSocket.Receive(req)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (s *Server) StartApp(jSocket *socket.JsonSocket, appKey string) {
	defer jSocket.Close()

	appConfig, exist := s.appsConfigs[appKey]
	if !exist {
		return
	}

	appServer, exist := s.appsServers[appKey]
	if exist {
		appServer.Close()
		delete(s.appsServers, appKey)
	}

	newAppServer := app.NewAppServer(appConfig.Port, jSocket)
	s.appsServers[appKey] = newAppServer

	defer newAppServer.Close()
	newAppServer.Listen()
}

func (s *Server) RegisterAppClientBackend(jSocket *socket.JsonSocket, appKey string, appClientId string) error {
	appServer, exist := s.appsServers[appKey]
	if !exist {
		return errors.New("NOTFOUND")
	}

	appServer.HandleAppClientBackend(appClientId, jSocket.Conn())

	return nil
}

func (s *Server) Close() {
	// TODO
}
