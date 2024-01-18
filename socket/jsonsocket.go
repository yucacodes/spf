package socket

import (
	"encoding/json"
	"net"
)

var transferEnd = byte(0x0a)

type JsonSocket struct {
	esocket *ESocket
}

func NewJsonSocket(conn net.Conn) *JsonSocket {
	eSocket := NewESocket(conn)
	o := JsonSocket{
		esocket: eSocket,
	}
	return &o
}

func (js *JsonSocket) Send(v any) error {
	out, err := json.Marshal(v)
	if err != nil {
		return err
	}
	err = js.esocket.Send(out, transferEnd)
	if err != nil {
		return err
	}
	return nil
}

func (js *JsonSocket) Receive(v any) error {
	buff := make([]byte, 0, 1000)
	buff, err := js.esocket.Receive(buff, transferEnd)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buff, v)
	if err != nil {
		return err
	}
	return nil
}

func (js *JsonSocket) IsOpen() bool {
	return js.esocket.IsOpen()
}

func (js *JsonSocket) IsClosed() bool {
	return js.esocket.IsClosed()
}
