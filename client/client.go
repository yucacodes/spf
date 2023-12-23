package client

import (
	"fmt"
	"net"

	"github.com/yucacodes/secure-port-forwarding/stream"
	"github.com/yucacodes/secure-port-forwarding/transfers"
)

func Main() {
	fmt.Println("Client")
	// Connect to the server
	server, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer server.Close()

	for {
		err = transfers.Write(server, 0)
		if err != nil {
			fmt.Println("Sending start code Error")
			continue
		}
		break
	}

	fmt.Println("Star code 0 sent")

	for {
		var code int
		err := transfers.Read(server, &code)
		if err != nil {
			// fmt.Println("Error waiting for sub client code")
			continue
		}
		go handleCallback(server, code)
	}
}

func handleCallback(server net.Conn, startCode int) {
	callback, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting callback:", err)
		return
	}
	fmt.Println("Sending callback code", startCode)
	err = transfers.Write(callback, startCode)
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
