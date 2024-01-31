package service

import (
	"fmt"
	"log"
	"net"

	"github.com/yucacodes/secure-port-forwarding/config"
	"github.com/yucacodes/secure-port-forwarding/socket"
	"golang.org/x/sync/syncmap"
)

type ForeignService struct {
	ownerConnection    *socket.JsonSocket
	config             config.Service
	clientsConnections *syncmap.Map
	logger             *log.Logger
}

type ForeignServiceClientConectionPairRequest struct {
	ClientId string
}

func (s *ForeignService) HandleIncomingClientConnection(clientId string, conn net.Conn) {
	clientConn := NewConnectionPair(clientId, conn)
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
