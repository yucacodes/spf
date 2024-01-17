package testclient

import (
	"fmt"
	"net"
	"time"

	"github.com/yucacodes/secure-port-forwarding/sockets"
)

func Main() {
	fmt.Println("Client")

	conn, err := net.Dial("tcp", "localhost:5000")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()

	eSocket := sockets.NewESocket(conn)

	for i := 0; i < 5 && !eSocket.IsClosed(); i++ {
		eSocket.SendWithStop([]byte{1, 2, 3, 4}, 0)
		time.Sleep(1000 * time.Millisecond)
	}
}
