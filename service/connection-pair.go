package service

import (
	"net"

	"github.com/google/uuid"
	"github.com/yucacodes/secure-port-forwarding/socket"
)

type ConnectionPair struct {
	clientId           string
	incomingConnection net.Conn
	backendConnection  net.Conn
}

func NewConnectionPair(incomingConnection net.Conn) *ConnectionPair {
	return &ConnectionPair{
		clientId:           uuid.New().String(),
		incomingConnection: incomingConnection,
	}
}

func (c *ConnectionPair) ClientId() string {
	return c.clientId
}

func (c *ConnectionPair) SetBackend(backendConnection net.Conn) {
	c.backendConnection = backendConnection
}

func (c *ConnectionPair) Streaming() {
	socket.NewESocket(c.incomingConnection).PairStreaming(socket.NewESocket(c.backendConnection))
}