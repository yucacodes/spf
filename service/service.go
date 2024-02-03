package service

import "net"

type Service interface {
	HandleIncomingClientConnection(conn net.Conn, clientId *string)
	HandleBackendServiceConnection(clientId string, conn net.Conn)
	Stop()
}
