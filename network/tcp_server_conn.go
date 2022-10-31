package network

import (
	"net"

	uuid "github.com/satori/go.uuid"
)

type ConnID = uuid.UUID

type TCPServerConn struct {
	ID   ConnID
	conn net.Conn // tcp conn
}

func (c TCPServerConn) Send(message []byte) error {
	_, err := c.conn.Write(message)
	return err
}

func (c TCPServerConn) close() error {
	err := c.conn.Close()
	return err
}
