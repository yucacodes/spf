package app

import (
	"net"
	"strconv"

	"github.com/yucacodes/secure-port-forwarding/socket"
)

type AppServer struct {
	clients map[string]*AppClient
	backend *socket.JsonSocket
	port    int
}

func NewAppServer(port int, backend *socket.JsonSocket) *AppServer {
	o := AppServer{port: port, backend: backend}
	return &o
}

func (as *AppServer) Listen() error {

	defer as.Close()

	server, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(as.port))
	if err != nil {
		return err
	}

	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			return err
		}

		client := NewAppClient(conn)
		as.clients[client.Id()] = client
		go as.RequestAppClientBackend(client)
	}

}

func (as *AppServer) HandleAppClientBackend(clientId string, conn net.Conn) {
	client, exist := as.clients[clientId]
	if !exist {
		conn.Close()
		return
	}
	client.SetBackendConnection(conn)
	go client.Streaming()
}

type AppClientPairRequestDto struct {
	ClientId string
}

func (as *AppServer) RequestAppClientBackend(client *AppClient) error {
	dto := AppClientPairRequestDto{ClientId: client.Id()}
	err := as.backend.Send(dto)
	if err != nil {
		return err
	}
	return nil
}

func (as *AppServer) Close() {
	as.backend.Close()
	for clientId, client := range as.clients {
		client.Close()
		delete(as.clients, clientId)
	}
}
