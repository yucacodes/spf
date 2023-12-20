package stream

import (
	"fmt"
	"net"
)

func repeat(from net.Conn, to net.Conn) {
	buf := make([]byte, 1)
	for {
		n, _ := from.Read(buf)
		if n > 0 {
			to.Write(buf)
			fmt.Print(string(buf))
		}
	}
}

func HandlePairStream(a net.Conn, b net.Conn) {
	fmt.Println("HandlePairStream")
	go repeat(a, b)
	fmt.Println("A > B")
	go repeat(b, a)
	fmt.Println("B > A")
	for {
	}
}
