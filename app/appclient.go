package app

import (
	"net"

	"github.com/google/uuid"
)

type AppClient struct {
	id         string
	clientConn net.Conn
	backConn   net.Conn
}

func NewAppClient(clientConn net.Conn) *AppClient {
	o := AppClient{
		id:         uuid.New().String(),
		clientConn: clientConn,
	}
	return &o
}

func (ap *AppClient) Id() string {
	return ap.id
}

func (ap *AppClient) SetBackendConnection(backConn net.Conn) {
	ap.backConn = backConn
}

func (ap *AppClient) Streaming() {
	// TODO: Implement this
	panic("Implement this")
}

func (ap *AppClient) Close() {
	if ap.backConn != nil {
		ap.backConn.Close()
	}
	if ap.clientConn != nil {
		ap.clientConn.Close()
	}
}
