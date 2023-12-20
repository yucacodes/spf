package server

import (
	"fmt"
	"net"

	"github.com/yucacodes/secure-port-forwarding/stream"
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

	subClientsPool := make(map[int]*net.Conn)

	for {
		// Accept incoming connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// Handle client connection in a goroutine
		fmt.Println("Client connected")
		go handleClient(conn, subClientsPool)
	}
}

func handleClient(conn net.Conn, subClientsPool map[int]*net.Conn) {
	defer conn.Close()
	var startCode int

	err := transfers.Receive(conn, &startCode)
	if err != nil {
		fmt.Print("Error reading Start Code")
		return
	}
	fmt.Println("Start Code received", startCode)

	if startCode == 0 {
		handleMainClient(conn, subClientsPool)
	} else {
		subClient, exist := subClientsPool[startCode]
		if !exist || subClient == nil {
			return
		}
		stream.HandlePairStream(conn, *subClient)
	}
}

func handleMainClient(client net.Conn, pool map[int]*net.Conn) {
	defer client.Close()
	listener, err := net.Listen("tcp", "0.0.0.0:9000")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server is listening on port 9000")

	for {
		// Accept incoming connections
		subClient, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		fmt.Println("Sub Client connected", subClient)
		subClientCode := 500
		err = transfers.Send(client, subClientCode)
		if err != nil {
			fmt.Println("Error on transmit code to subclient")
		}
		pool[subClientCode] = &subClient
	}
}
