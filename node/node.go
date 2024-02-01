package node

import (
	"errors"
	"log"
	"net"

	"github.com/yucacodes/secure-port-forwarding/config"
	"github.com/yucacodes/secure-port-forwarding/listen"
	"github.com/yucacodes/secure-port-forwarding/request"
	"github.com/yucacodes/secure-port-forwarding/service"
	"github.com/yucacodes/secure-port-forwarding/socket"
	"golang.org/x/sync/syncmap"
)

type Node struct {
	config            *config.Config
	availableServices *syncmap.Map
	logger            *log.Logger
}

func (node *Node) Run() {
	for _, listenConfig := range node.config.Listen {
		listen := listen.NewListen(&listenConfig, node.availableServices)
		go listen.Start()
		defer listen.Stop()
	}

	// for _, publishConfig := range node.config.Publish {
	// TODO
	// 	publish := publish.NewPublishedServiceThroughNode(node.config.Id, )
	// }

	node.ListenNodeRequests()
}

func (node *Node) ListenNodeRequests() {
	server, err := net.Listen("tcp", node.config.ListenConnection())
	if err != nil {
		node.logger.Println(err)
		return
	}
	defer server.Close()

	node.logger.Println("Listening on " + node.config.ListenConnection())
	for {
		conn, err := server.Accept()
		if err != nil {
			node.logger.Println(err)
			break
		}
		go func() {
			node.HandleConnection(conn)
			conn.Close()
		}()
	}
}

func (node *Node) HandleConnection(conn net.Conn) {
	node.logger.Println("New connection")
	jSocket := socket.NewJsonSocket(conn)
	req, err := node.GetNodeRequest(jSocket)
	if err != nil {
		node.logger.Println("Error reading App request")
		return
	}

	if req.PublishService != nil {
		node.logger.Println("Received Publish Service request")
		node.CreateForeignService(conn, req.PublishService)
	} else if req.StreamingToServiceClient != nil {
		node.logger.Println("Received Streaming Service Client request")
		node.StreamToServiceClient(conn, req.StreamingToServiceClient)
	}
}

func (node *Node) GetNodeRequest(jSocket *socket.JsonSocket) (*request.NodeRequest, error) {
	req := request.NodeRequest{}
	err := jSocket.Receive(&req)
	if err != nil {
		return nil, err
	}
	if node.config.DisableNodeValidation != nil && *node.config.DisableNodeValidation {
		return &req, nil
	}
	for _, nodeConfig := range node.config.Nodes {
		if req.Id.Name == nodeConfig.Name && nodeConfig.Key != nil && req.Id.Key == *nodeConfig.Key {
			return &req, nil
		}
	}
	return nil, errors.New("node not found")
}

func (node *Node) CreateForeignService(conn net.Conn, req *request.PublishServiceRequest) {
	_oldService, exist := node.availableServices.Load(req.Service)
	if exist {
		oldService := _oldService.(service.Service)
		oldService.Stop()
	}
	newService := service.NewForeignService(conn)
	node.availableServices.Store(req.Service, newService)
}

func (node *Node) StreamToServiceClient(conn net.Conn, req *request.StreamingToServiceClient) {
	_service, exist := node.availableServices.Load(req.Service)

	if !exist {
		node.logger.Println("Not found requested app server")
		return
	}
	service := _service.(service.Service)
	service.HandleBackendServiceConnection(req.Client, conn)
}
