package app

import (
	"log"
	"net"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/yucacodes/secure-port-forwarding/socket"
)

type AppClient struct {
	id         string
	clientConn net.Conn
	backConn   net.Conn
	logger     *log.Logger
}

func NewAppClient(clientConn net.Conn) *AppClient {
	o := AppClient{
		id:         uuid.New().String(),
		clientConn: clientConn,
		logger:     log.New(os.Stdout, "AppClient: ", log.Ldate|log.Ltime),
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
	ap.logger.Println("Starting streaming...")
	cSocket := socket.NewESocket(ap.clientConn)
	bSocket := socket.NewESocket(ap.backConn)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		cSocket.StreamingTo(bSocket)
		bSocket.Conn().Close()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		bSocket.StreamingTo(cSocket)
		cSocket.Conn().Close()
	}()

	wg.Wait()
	ap.logger.Println("Streaming Finished")
}
