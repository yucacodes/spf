package node

import (
	"log"
	"net"

	"github.com/yucacodes/secure-port-forwarding/config"
	"github.com/yucacodes/secure-port-forwarding/publish"
	"github.com/yucacodes/secure-port-forwarding/request"
	"github.com/yucacodes/secure-port-forwarding/socket"
)

type Node struct {
	config            *config.Config
	logger            *log.Logger
	publishedServices map[string]publish.PublishedService
}

func (node *Node) Run() {
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
		node.PublishForeignService(conn, req.PublishService)
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
	return &req, nil
}

func (node *Node) PublishForeignService(conn net.Conn, req *request.PublishServiceRequest) {
	ps := publish.NewPublishedForeignService()
	node.publishedServices[req.Service] = ps
	ps.Start()
}

func (node *Node) StreamToServiceClient(conn net.Conn, req *request.StreamingToServiceClient) {
	publishedService, exist := node.publishedServices[req.Service]
	if !exist {
		node.logger.Println("Not found requested app server")
		return
	}
	publishedService.HandleServiceClientBackend(req.Client, conn)
}
