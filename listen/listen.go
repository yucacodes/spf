package listen

import (
	"log"
	"net"
	"os"

	"github.com/yucacodes/secure-port-forwarding/config"
	"github.com/yucacodes/secure-port-forwarding/service"
	"golang.org/x/sync/syncmap"
)

type Listen struct {
	config            *config.Listen
	availableServices *syncmap.Map
	server            net.Listener
	logger            *log.Logger
}

func NewListen(config *config.Listen, availableServices *syncmap.Map) *Listen {
	return &Listen{
		config:            config,
		availableServices: availableServices,
		logger:            log.New(os.Stdout, "Listen: ", log.Ldate|log.Ltime),
	}
}

func (l *Listen) Start() {
	server, err := net.Listen("tcp", l.config.ListenConnection())
	if err != nil {
		l.logger.Println(err)
		return
	}
	l.logger.Println("Listening on " + l.config.ListenConnection())
	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			break
		}
		l.logger.Println("New connection")

		go func() {
			_service, exist := l.availableServices.Load(l.config.Connect.Service)
			if exist {
				service := _service.(service.Service)
				service.HandleIncomingClientConnection(conn)
			} else {
				conn.Close()
			}
		}()

	}
	l.logger.Println("Releasing " + l.config.ListenConnection())
}

func (l *Listen) Stop() {
	l.server.Close()
}
