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
	conn, err := net.Dial("tcp", c.serverHost+":"+strconv.Itoa(c.serverPort))
	if err != nil {
		return
	}
	c.server = socket.NewJsonSocket(conn)
	defer c.server.Close()

	req := server.AppRequest{
		AppKey:  c.appKey,
		InitApp: true,
	}
	err = c.server.Send(req)
	if err != nil {
		c.logger.Fatalln(err)
		return
	}

	for c.server.IsOpen() {
		req := app.AppClientPairRequestDto{}
		err := c.server.Receive(req)
		if err != nil {
			continue
		}

		go c.createBackendConnection(&req)
	}
}

func (c *Client) createBackendConnection(req *app.AppClientPairRequestDto) {

	backConn, err := net.Dial("tcp", c.appHost+":"+strconv.Itoa(c.appPort))
	if err != nil {
		return
	}
	defer backConn.Close()

	serverConn, err := net.Dial("tcp", c.serverHost+":"+strconv.Itoa(c.serverPort))
	if err != nil {
		return
	}
	sjSocket := socket.NewJsonSocket(serverConn)
	defer sjSocket.Close()

	serverReq := server.AppRequest{
		AppKey:             c.appKey,
		BackendToAppClient: true,
		AppClientId:        req.ClientId,
	}
	err = sjSocket.Send(serverReq)
	if err != nil {
		return
	}

	appClient := app.NewAppClient(serverConn)
	appClient.SetBackendConnection(backConn)
	appClient.Streaming()
}

func (c *Client) Close() {
	// TODO
}
