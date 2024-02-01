package service

import "net"

type Service interface {
	HandleIncomingClientConnection(conn net.Conn)
	HandleBackendServiceConnection(clientId string, conn net.Conn)
	Stop()
}
