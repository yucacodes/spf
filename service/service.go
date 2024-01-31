package service

import "net"

type Service interface {
	HandleIncomingClientConnection(clientId string, conn net.Conn)
	HandleBackendServiceConnection(clientId string, conn net.Conn)
}
