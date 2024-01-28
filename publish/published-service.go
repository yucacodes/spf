package publish

import "net"

type PublishedService interface {
	HandleServiceClientBackend(string, net.Conn)
}
