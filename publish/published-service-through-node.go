package publish

import (
	"log"
	"net"
	"os"

	"github.com/yucacodes/secure-port-forwarding/config"
	"github.com/yucacodes/secure-port-forwarding/request"
	"github.com/yucacodes/secure-port-forwarding/service"

	"github.com/yucacodes/secure-port-forwarding/socket"
)

type PublishedServiceThroughNode struct {
	id         *config.NodeId
	service    *config.Service
	nodeConfig *config.Node
	nodeConn   *socket.JsonSocket
	logger     *log.Logger
}

func NewPublishedServiceThroughNode(
	id *config.NodeId,
	service *config.Service,
	node *config.Node,
) *PublishedServiceThroughNode {
	ps := PublishedServiceThroughNode{
		id:         id,
		service:    service,
		nodeConfig: node,
		logger:     log.New(os.Stdout, "Client: ", log.Ldate|log.Ltime),
	}
	return &ps
}

func (c *PublishedServiceThroughNode) Connect() {
	if !c.nodeConfig.IsPublic() {
		c.logger.Println("Error: Node " + c.nodeConfig.Name + " is not public")
		return
	}
	if !c.service.IsOwn() {
		c.logger.Println("Error: Service " + c.service.Name + " is not own")
		return
	}
	c.logger.Println("Connecting to the node " + c.nodeConfig.Name)
	conn, err := net.Dial("tcp", c.nodeConfig.Connection())
	if err != nil {
		c.logger.Println("Connection error")
		c.logger.Println(err)
		return
	}
	defer conn.Close()

	c.nodeConn = socket.NewJsonSocket(conn)

	c.logger.Println("Sending publish service request")
	req := request.NodeRequest{
		Id: *c.id,
		PublishService: &request.PublishServiceRequest{
			Service: c.service.Name,
		},
	}
	err = c.nodeConn.Send(req)
	if err != nil {
		c.logger.Println("Error sending publish service request")
		c.logger.Println(err)
		return
	}

	c.logger.Println("Waiting for connection requests...")
	for c.nodeConn.IsOpen() {
		req := service.ForeignServiceClientConectionPairRequest{}
		err := c.nodeConn.Receive(&req)
		if err != nil {
			c.logger.Println("Error reading connection request")
			continue
		}
		c.logger.Println("Connection request received")

		go c.createBackendConnection(&req)
	}
}

func (c *PublishedServiceThroughNode) createBackendConnection(req *service.ForeignServiceClientConectionPairRequest) {
	c.logger.Println("Connecting to backend...")
	serviceConn, err := net.Dial("tcp", c.service.Connection())

	if err != nil {
		c.logger.Println("connection error")
		c.logger.Println(err)
		return
	}
	defer serviceConn.Close()

	c.logger.Println("Connecting to node...")
	serverConn, err := net.Dial("tcp", c.nodeConfig.Connection())
	if err != nil {
		c.logger.Println("connection error")
		c.logger.Println(err)
		return
	}
	defer serverConn.Close()

	sjSocket := socket.NewJsonSocket(serverConn)

	serverReq := request.NodeRequest{
		StreamingToServiceClient: &request.StreamingToServiceClient{
			Service: c.service.Name,
			Client:  req.ClientId,
		},
	}

	c.logger.Println("Sending request to server for stablish service client streaming...")
	err = sjSocket.Send(serverReq)
	if err != nil {
		c.logger.Println("Sending request error")
		c.logger.Println(err)
		return
	}

	connectionPair := service.NewConnectionPair(serverConn)
	connectionPair.SetBackend(serviceConn)
	c.logger.Println("Starting App client streaming")
	connectionPair.Streaming()
}
