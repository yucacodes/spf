package server

import (
	"fmt"
	"net"

	"github.com/yucacodes/secure-port-forwarding/transfers"
)

func Main() {
	fmt.Println("Server")
	listener, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server is listening on port 8080")

	for {
		// Accept incoming connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// Handle client connection in a goroutine
		fmt.Println("Client connected")
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	var receive int
	err := transfers.Receive(conn, &receive)
	if err != nil {
		fmt.Printf("Receiving Error")
		fmt.Println(err)
	}
	fmt.Println(receive)
}
