package network

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

type TCPClient interface {
	Close()
	OnClose() <-chan bool
	Dia(url string) error
	Send(msg string) error
}

type Client struct {
	conn        net.Conn
	receive     chan []byte
	writer      io.Writer
	closeSignal chan bool
}

// OnClose implements TCPClient
func (c *Client) OnClose() <-chan bool {
	return c.closeSignal
}

func NewTCPClient(writer io.Writer) TCPClient {
	client := &Client{
		receive:     make(chan []byte),
		closeSignal: make(chan bool),
		writer:      writer,
	}
	return client
}

func (c *Client) Dia(url string) error {

	// connect to this socket
	conn, err := net.Dial("tcp", url)
	if err != nil {
		fmt.Fprintln(c.writer, "Could not dial url: "+url)
		return err
	}
	c.conn = conn
	defer c.Close()
	for {
		message, err := bufio.NewReader(c.conn).ReadBytes('\n')
		if err != nil {
			return err
		}
		fmt.Fprintf(c.writer, "client pump receive: "+string(message))
	}
}

func (c *Client) Send(msg string) error {
	// NOTE: need to add '\n' as line end
	_, err := c.conn.Write([]byte(msg + "\n"))
	return err
}

// OnClose implements TCPClient
func (c *Client) Close() {
	c.conn.Close()
	c.closeSignal <- true
}
