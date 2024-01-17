package testserver

import (
	"fmt"
	"net"

	"github.com/yucacodes/secure-port-forwarding/sockets"
)

func Main() {
	server, err := net.Listen("tcp", "0.0.0.0:5000")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer server.Close()
	fmt.Println("Server is listening")

	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		fmt.Println("New connection")
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	eSocket := sockets.NewESocket(conn)
	defer eSocket.Close()

	for eSocket.IsOpen() {
		read, err := eSocket.ReceiveUntilStop(0, []byte{})
		if err == nil {
			fmt.Println(read)
		}
	}

	fmt.Println("End connection")
}
