package client

import (
	"fmt"
	"net"

	"github.com/yucacodes/secure-port-forwarding/transfers"
)

func Main() {
	fmt.Println("Client")
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()

	err = transfers.Send(conn, 4500)
	if err != nil {
		fmt.Printf("Sending Error")
	}
}
