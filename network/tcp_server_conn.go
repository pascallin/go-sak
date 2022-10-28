package network

import (
	"fmt"
	"net"

	uuid "github.com/satori/go.uuid"
)

type ConnID = uuid.UUID

type TCPServerConn struct {
	ID   ConnID
	conn net.Conn // tcp conn
}

func (c *TCPServerConn) send(message string) error {
	_, err := fmt.Fprintln(c.conn, message)
	return err
}

func (c *TCPServerConn) close() error {
	err := c.conn.Close()
	return err
}
