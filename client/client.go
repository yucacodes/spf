package client

import (
	"fmt"
	"net"
	"strconv"

	"github.com/yucacodes/secure-port-forwarding/app"
	"github.com/yucacodes/secure-port-forwarding/server"
	"github.com/yucacodes/secure-port-forwarding/socket"
	"github.com/yucacodes/secure-port-forwarding/stream"
)

type Client struct {
	host   string
	port   int
	appKey string
	server *socket.JsonSocket
}

func (c *Client) Start() {
	conn, err := net.Dial("tcp", c.host+":"+strconv.Itoa(c.port))
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
		return
	}

	for c.server.IsOpen() {
		req := app.AppClientPairRequestDto{}
		err := c.server.Receive(req)
		if err != nil {
			continue
		}

	}
}

func handleCallback(server net.Conn, appClientId string) {
	callback, err := net.Dial("tcp", "localhost:5000")
	if err != nil {
		fmt.Println("Error connecting callback:", err)
		return
	}
	fmt.Println("Sending callback code", appClientId)

	err = transfers.Send(callback, appClientId)
	if err != nil {
		fmt.Println("Error transfering callback code")
		fmt.Println(err)
		return
	}
	backend, err := net.Dial("tcp", "localhost:5173")
	if err != nil {
		fmt.Println("Error connecting backend:", err)
		return
	}
	stream.HandlePairStream(backend, callback)
}
