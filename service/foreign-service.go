package service

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/yucacodes/secure-port-forwarding/socket"
	"golang.org/x/sync/syncmap"
)

type ForeignService struct {
	ownerConnection    *socket.JsonSocket
	clientsConnections *syncmap.Map
	logger             *log.Logger
}

func NewForeignService(ownerConnection net.Conn) *ForeignService {
	return &ForeignService{
		ownerConnection:    socket.NewJsonSocket(ownerConnection),
		clientsConnections: &syncmap.Map{},
		logger:             log.New(os.Stdout, "ForeignService: ", log.Ldate|log.Ltime),
	}
}

type ForeignServiceClientConectionPairRequest struct {
	ClientId string
}

func (s *ForeignService) HandleIncomingClientConnection(conn net.Conn) {
	clientConn := NewConnectionPair(conn)
	s.clientsConnections.Store(clientConn.ClientId(), clientConn)

	req := ForeignServiceClientConectionPairRequest{ClientId: clientConn.ClientId()}
	err := s.ownerConnection.Send(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	s.logger.Println("Request App client backend success")
}

func (s *ForeignService) HandleBackendServiceConnection(clientId string, conn net.Conn) {
	_clientConn, exist := s.clientsConnections.Load(clientId)
	if !exist {
		s.logger.Println("Not found requested app client")
		return
	}
	clientConn := _clientConn.(*ConnectionPair)

	clientConn.SetBackend(conn)
	clientConn.Streaming()
}

func (s *ForeignService) Stop() {
	// s.clientsConnections close
	s.ownerConnection.Conn().Close()
}
