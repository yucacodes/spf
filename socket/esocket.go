package socket

import (
	"errors"
	"fmt"
	"net"
)

type ESocket struct {
	conn   net.Conn
	closed bool
}

func NewESocket(conn net.Conn) *ESocket {
	es := ESocket{
		conn:   conn,
		closed: false,
	}
	return &es
}

func (es *ESocket) SendWithStop(buff []byte, stopByte byte) error {
	if es.IsClosed() {
		return errors.New("EOF")
	}
	out := append(buff, stopByte)
	n, err := es.conn.Write(out)
	if n < len(out) {
		es.closed = true
		return errors.New("EOF")
	}
	if err != nil {
		return err
	}
	return nil
}

func (es *ESocket) StreamingTo(target *ESocket) {
	for es.IsOpen() && target.IsOpen() {
		b, err := es.ReceiveByte()
		if err != nil {
			fmt.Println(err)
			continue
		}
		target.SendByte(b)
	}
}

func (es *ESocket) SendByte(v byte) error {
	if es.IsClosed() {
		return errors.New("EOF")
	}
	out := []byte{v}
	n, err := es.conn.Write(out)
	if n < 1 {
		es.closed = true
		return errors.New("EOF")
	}
	if err != nil {
		return err
	}
	return nil
}

func (es *ESocket) ReceiveByte() (byte, error) {
	if es.IsClosed() {
		return 0, errors.New("EOF")
	}
	oneByteBuff := make([]byte, 1)
	n, err := es.conn.Read(oneByteBuff)
	if err != nil || n == 0 {
		es.closed = true
		return 0, errors.New("EOF")
	}
	return oneByteBuff[0], nil
}

func (es *ESocket) WaitForClose() {
	for es.IsOpen() {
		es.ReceiveByte()
	}
}

func (es *ESocket) ReceiveUntilStop(buff []byte, stopByte byte) ([]byte, error) {
	if es.IsClosed() {
		return nil, errors.New("EOF")
	}
	mainBuf := buff
	oneByteBuff := make([]byte, 1)
	for {
		n, err := es.conn.Read(oneByteBuff)
		if err != nil || n == 0 {
			es.closed = true
			return nil, errors.New("EOF")
		}
		if oneByteBuff[0] == stopByte {
			break
		}
		mainBuf = append(mainBuf, oneByteBuff...)
	}
	return mainBuf, nil
}

func (es *ESocket) IsClosed() bool {
	return es.closed
}

func (es *ESocket) IsOpen() bool {
	return !es.closed
}

func (es *ESocket) Conn() net.Conn {
	return es.conn
}
