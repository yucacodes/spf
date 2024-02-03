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
	id                        *config.NodeId
	name                      string
	directProvider            *config.Node
	reverseProviderConnection *socket.JsonSocket
	clientsConnections        *syncmap.Map
	logger                    *log.Logger
}

func NewForeignService(id *config.NodeId, name string, directProvider *config.Node, reverseProviderConnection net.Conn) *ForeignService {
	var _reverseProviderConnection *socket.JsonSocket = nil
	if reverseProviderConnection != nil {
		_reverseProviderConnection = socket.NewJsonSocket(reverseProviderConnection)
	}
	fs := &ForeignService{
		id:                        id,
		name:                      name,
		reverseProviderConnection: _reverseProviderConnection,
		directProvider:            directProvider,
		clientsConnections:        &syncmap.Map{},
		logger:                    log.New(os.Stdout, "ForeignService: ", log.Ldate|log.Ltime),
	}

	if fs.directProvider == nil && fs.reverseProviderConnection == nil {
		fs.logger.Println("Error: none providers")
		return nil
	}

	if fs.directProvider != nil && fs.reverseProviderConnection != nil {
		fs.logger.Println("Error: multiple type of providers")
		return nil
	}

	if fs.directProvider != nil && !fs.directProvider.IsPublic() {
		fs.logger.Println("Error: Node " + fs.directProvider.Name + " is not public")
		return nil
	}

	return fs
}

func (s *ForeignService) HandleIncomingClientConnection(conn net.Conn, clientId *string) {
	if s.directProvider != nil {
		s.handleIncomingClientConnectionWithDirectProvider(conn, clientId)
	} else if s.reverseProviderConnection != nil {
		s.handleIncomingClientConnectionWithReverseProvider(conn, clientId)
	}
}

func (s *ForeignService) handleIncomingClientConnectionWithReverseProvider(conn net.Conn, clientId *string) {
	clientConn := NewConnectionPair(conn, clientId)
	s.clientsConnections.Store(clientConn.ClientId(), clientConn)
	req := request.NodeRequest{
		Id:                                *s.id,
		ForeignServiceClientConectionPair: &request.ForeignServiceClientConectionPairRequest{Client: clientConn.ClientId()},
	}
	err := s.reverseProviderConnection.Send(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	s.logger.Println("Request Service connection backend success")
}

func (s *ForeignService) handleIncomingClientConnectionWithDirectProvider(conn net.Conn, clientId *string) {
	clientConn := NewConnectionPair(conn, clientId)
	s.clientsConnections.Store(clientConn.ClientId(), clientConn)

	s.logger.Println("Connecting to the node " + s.directProvider.Name + " (" + s.directProvider.Connection() + ")")
	providerConn, err := net.Dial("tcp", s.directProvider.Connection())
	if err != nil {
		s.logger.Println("Connection error")
		s.logger.Println(err)
		return
	}
	defer providerConn.Close()

	providerJConn := socket.NewJsonSocket(providerConn)

	req := request.NodeRequest{
		Id: *s.id,
		ForeignServiceClientConectionPair: &request.ForeignServiceClientConectionPairRequest{
			Client:  clientConn.ClientId(),
			Service: s.name,
		},
	}

	err = providerJConn.Send(req)
	if err != nil {
		s.logger.Println("Error sending foreign service client conection pair request")
		s.logger.Println(err)
		return
	}

	clientConn.SetBackend(providerConn)
	clientConn.Streaming()
}

func (s *ForeignService) HandleBackendServiceConnection(clientId string, conn net.Conn) {
	_clientConn, exist := s.clientsConnections.Load(clientId)
	if !exist {
		s.logger.Println("Not found requested service client")
		return
	}
	clientConn := _clientConn.(*ConnectionPair)

	clientConn.SetBackend(conn)
	clientConn.Streaming()
}

func (s *ForeignService) Stop() {
	// TODO: s.clientsConnections close ?
	s.reverseProviderConnection.Conn().Close()
}
