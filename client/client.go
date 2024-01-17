package client

import (
	"fmt"
	"net"

	transfers "github.com/yucacodes/secure-port-forwarding/messages"
	"github.com/yucacodes/secure-port-forwarding/server"
	"github.com/yucacodes/secure-port-forwarding/stream"
)

func Main() {
	fmt.Println("Client")
	// Connect to the serverListener
	serverListener, err := net.Dial("tcp", "localhost:5000")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer serverListener.Close()

	for {
		appRequest := server.AppRequest{
			Key:              "alsdknaslkdnasodnalsdknasdoasdasdasdasdlkamdsasd",
			SetAppConnection: true,
		}

		err = transfers.Send(serverListener, appRequest)
		if err != nil {
			fmt.Println("Sending start code Error")
			continue
		}
		break
	}

	fmt.Println("Star code 0 sent")

	for {
		connectClientRequest := server.ConnectClientRequest{}
		err := transfers.Receive(serverListener, &connectClientRequest)
		if err != nil {
			// fmt.Println("Error waiting for sub client code")
			continue
		}
		go handleCallback(serverListener, connectClientRequest.ClientId)
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
