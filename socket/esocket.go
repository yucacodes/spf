package socket

import (
	"errors"
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

func (es *ESocket) Send(buff []byte, stopByte byte) error {
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

func (es *ESocket) Receive(buff []byte, stopByte byte) ([]byte, error) {
	if es.IsClosed() {
		return nil, errors.New("EOF")
	}
	mainBuf := buff
	oneByteBuff := make([]byte, 1)
	for {
		_, err := es.conn.Read(oneByteBuff)
		if err != nil {
			if err.Error() == "EOF" {
				es.closed = true
			}
			return nil, err
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
