package app

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/yucacodes/secure-port-forwarding/socket"
)

type AppServer struct {
	clients map[string]*AppClient
	backend *socket.JsonSocket
	port    int
	logger  *log.Logger
}

func NewAppServer(port int, backend net.Conn) *AppServer {
	o := AppServer{
		clients: make(map[string]*AppClient),
		port:    port,
		backend: socket.NewJsonSocket(backend),
		logger:  log.New(os.Stdout, "AppServer: ", log.Ldate|log.Ltime),
	}
	return &o
}

func (as *AppServer) Listen() {

	server, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(as.port))
	if err != nil {
		as.logger.Println(err)
		return
	}
	as.logger.Println("Listening on port " + strconv.Itoa(as.port))
	defer server.Close()

	go func() {
		as.backend.WaitForClose()
		server.Close()
	}()

	for {
		conn, err := server.Accept()
		if err != nil {
			break
		}
		as.logger.Println("New connection")
		client := NewAppClient(conn)
		as.clients[client.Id()] = client
		go as.RequestAppClientBackend(client)
	}

	as.logger.Println("Releasing port " + strconv.Itoa(as.port))
}

func (as *AppServer) HandleAppClientBackend(clientId string, conn net.Conn) {
	client, exist := as.clients[clientId]
	if !exist {
		as.logger.Println("Not found requested app client")
		return
	}
	client.SetBackendConnection(conn)
	client.Streaming()
}

type AppClientPairRequestDto struct {
	ClientId string
}

func (as *AppServer) RequestAppClientBackend(client *AppClient) {
	dto := AppClientPairRequestDto{ClientId: client.Id()}
	err := as.backend.Send(dto)
	if err != nil {
		fmt.Println(err)
		return
	}
	as.logger.Println("Request App client backend success")
}
