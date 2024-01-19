package client

import (
	"fmt"
	"net"
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
