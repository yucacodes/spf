package service

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/yucacodes/secure-port-forwarding/config"
	"github.com/yucacodes/secure-port-forwarding/request"
	"github.com/yucacodes/secure-port-forwarding/socket"
	"golang.org/x/sync/syncmap"
)

type ForeignService struct {
	id                 *config.NodeId
	providerConnection *socket.JsonSocket
	clientsConnections *syncmap.Map
	logger             *log.Logger
}

func NewForeignService(id *config.NodeId, ownerConnection net.Conn) *ForeignService {
	return &ForeignService{
		id:                 id,
		providerConnection: socket.NewJsonSocket(ownerConnection),
		clientsConnections: &syncmap.Map{},
		logger:             log.New(os.Stdout, "ForeignService: ", log.Ldate|log.Ltime),
	}
}

func (s *ForeignService) HandleIncomingClientConnection(conn net.Conn) {
	clientConn := NewConnectionPair(conn)
	s.clientsConnections.Store(clientConn.ClientId(), clientConn)
	req := request.NodeRequest{
		Id:                                *s.id,
		ForeignServiceClientConectionPair: &request.ForeignServiceClientConectionPairRequest{Client: clientConn.ClientId()},
	}
	err := s.providerConnection.Send(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	s.logger.Println("Request Service connection backend success")
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
	s.providerConnection.Conn().Close()
}
