package server

import (
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
		appsServers: make(map[string]*app.AppServer),
		logger:      log.New(os.Stdout, "Server: ", log.Ldate|log.Ltime),
	}

	for _, appConfig := range appsConfigs {
		s.appsConfigs[appConfig.Key] = appConfig
	}
	return &s
}

func (s *Server) Listen() {
	server, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(s.port))
	if err != nil {
		s.logger.Println(err)
		return
	}
	defer server.Close()

	s.logger.Println("Listening on port " + strconv.Itoa(s.port))
	for {
		conn, err := server.Accept()
		if err != nil {
			s.logger.Println(err)
			break
		}
		go func() {
			s.HandleConnection(conn)
			conn.Close()
		}()
	}
}

func (s *Server) HandleConnection(conn net.Conn) {
	s.logger.Println("New connection")
	jSocket := socket.NewJsonSocket(conn)
	req, err := s.GetAppRequest(jSocket)
	if err != nil {
		s.logger.Println("Error reading App request")
		return
	}

	if req.InitApp {
		s.logger.Println("Received Init App request")
		s.StartApp(conn, req.AppKey)
	} else if req.BackendToAppClient {
		s.logger.Println("Received Backend to App Client request")
		s.RegisterAppClientBackend(conn, req.AppKey, req.AppClientId)
	}
}

func (s *Server) GetNodeRequest(jSocket *socket.JsonSocket) (*NodeRequest, error) {
	req := NodeRequest{}
	err := jSocket.Receive(&req)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (s *Server) StartApp(conn net.Conn, appKey string) {
	appConfig, exist := s.appsConfigs[appKey]
	if !exist {
		s.logger.Println("Requested app not registered")
		return
	}

	_, exist = s.appsServers[appKey]
	if exist {
		s.logger.Println("Requested app was started previously, removing old app...")
		delete(s.appsServers, appKey)
	}

	newAppServer := app.NewAppServer(appConfig.Port, conn)
	s.appsServers[appKey] = newAppServer

	s.logger.Println("Starting app...")
	newAppServer.Listen()
}

func (s *Server) RegisterAppClientBackend(conn net.Conn, appKey string, appClientId string) {
	appServer, exist := s.appsServers[appKey]
	if !exist {
		s.logger.Println("Not found requested app server")
		return
	}

	appServer.HandleAppClientBackend(appClientId, conn)
}
