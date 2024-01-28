package publish

import (
	"net"

	"github.com/yucacodes/secure-port-forwarding/config"
)

type PublishedForeignService struct {
	service *config.Service
	node    *config.Node
}

func NewPublishedForeignService() *PublishedForeignService {
	return nil
}

func (pfs *PublishedForeignService) Start() {

}

func (pfs *PublishedForeignService) HandleServiceClientBackend(Client string, conn net.Conn) {

}
