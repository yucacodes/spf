package transfers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

var transferEnd = byte(0x0a)

func Send(conn net.Conn, v any) error {
	out, err := json.Marshal(v)
	if err != nil {
		return err
	}

	out = append(out, transferEnd)

	_, err = conn.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func Receive(conn net.Conn, v any) error {
	mainBuf := make([]byte, 0, 1000)
	charBuf := make([]byte, 1)

	for {
		n, err := conn.Read(charBuf)
		if n == 0 {
			if err != nil {
				return err
			} else {
				return errors.New("not receive data")
			}
		}
		if charBuf[len(charBuf)-1] == transferEnd {
			break
		}
		mainBuf = append(mainBuf, charBuf...)
	}
	err := json.Unmarshal(mainBuf, v)
	if err != nil {
		fmt.Println("Error parsing readed message")
		fmt.Print(err)
		return err
	}
	return nil
}
