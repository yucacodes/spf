package client

import (
	"log"
	"net"
	"os"
	"strconv"

	"github.com/yucacodes/secure-port-forwarding/app"
	"github.com/yucacodes/secure-port-forwarding/server"
	"github.com/yucacodes/secure-port-forwarding/socket"
)

type Client struct {
	serverHost string
	serverPort int
	appKey     string
	appHost    string
	appPort    int
	server     *socket.JsonSocket
	logger     *log.Logger
}

func NewClient(
	serverHost string,
	serverPort int,
	appKey string,
	appHost string,
	appPort int,
) *Client {
	c := Client{
		serverHost: serverHost,
		serverPort: serverPort,
		appKey:     appKey,
		appHost:    appHost,
		appPort:    appPort,
		logger:     log.New(os.Stdout, "Client: ", log.Ldate|log.Ltime),
	}
	return &c
}

func (c *Client) Connect() {
	c.logger.Println("Connecting to the server")
	conn, err := net.Dial("tcp", c.serverHost+":"+strconv.Itoa(c.serverPort))
	if err != nil {
		c.logger.Println("Connection error")
		c.logger.Println(err)
		return
	}
	defer conn.Close()

	c.server = socket.NewJsonSocket(conn)

	c.logger.Println("Sending Init App request")
	req := server.AppRequest{
		AppKey:  c.appKey,
		InitApp: true,
	}
	err = c.server.Send(req)
	if err != nil {
		c.logger.Println("Error sending init app request")
		c.logger.Println(err)
		return
	}

	c.logger.Println("Waiting for connection requests...")
	for c.server.IsOpen() {
		req := app.AppClientPairRequestDto{}
		err := c.server.Receive(&req)
		if err != nil {
			c.logger.Println("Error reading connection request")
			continue
		}
		c.logger.Println("Connection request received")

		go c.createBackendConnection(&req)
	}
}

func (c *Client) createBackendConnection(req *app.AppClientPairRequestDto) {
	c.logger.Println("Connecting to backend...")
	backConn, err := net.Dial("tcp", c.appHost+":"+strconv.Itoa(c.appPort))
	if err != nil {
		c.logger.Println("connection error")
		c.logger.Println(err)
		return
	}
	defer backConn.Close()

	c.logger.Println("Connecting to server...")
	serverConn, err := net.Dial("tcp", c.serverHost+":"+strconv.Itoa(c.serverPort))
	if err != nil {
		c.logger.Println("connection error")
		c.logger.Println(err)
		return
	}
	defer serverConn.Close()

	sjSocket := socket.NewJsonSocket(serverConn)

	serverReq := server.AppRequest{
		AppKey:             c.appKey,
		BackendToAppClient: true,
		AppClientId:        req.ClientId,
	}
	c.logger.Println("Sending request to server for set app client backend...")
	err = sjSocket.Send(serverReq)
	if err != nil {
		c.logger.Println("Sending request error")
		c.logger.Println(err)
		return
	}

	appClient := app.NewAppClient(serverConn)
	appClient.SetBackendConnection(backConn)
	c.logger.Println("Starting App client streaming")
	appClient.Streaming()
}
