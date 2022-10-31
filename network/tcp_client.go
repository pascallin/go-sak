package network

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

type TCPClient interface {
	Close()
	OnClose() <-chan bool
	Dia(url string) error
	Send(msg []byte) error
	OnReceive() <-chan []byte
	SetWriter(writer io.Writer)
}

type Client struct {
	conn        net.Conn
	receive     chan []byte
	writer      io.Writer
	closeSignal chan bool
}

// SetWriter implements TCPClient
func (c *Client) SetWriter(writer io.Writer) {
	c.writer = writer
}

// OnReceive implements TCPClient
func (c *Client) OnReceive() <-chan []byte {
	return c.receive
}

// OnClose implements TCPClient
func (c *Client) OnClose() <-chan bool {
	return c.closeSignal
}

func NewTCPClient() TCPClient {
	client := &Client{
		receive:     make(chan []byte),
		closeSignal: make(chan bool),
	}
	return client
}

func (c *Client) Dia(url string) error {
	// connect to this socket
	conn, err := net.Dial("tcp", url)
	if err != nil {
		if c.writer != nil {
			fmt.Fprintln(c.writer, "Could not dial url: "+url)
		}
		return err
	}
	c.conn = conn
	go c.onMessage()
	return nil
}

func (c *Client) Send(msg []byte) error {
	if c.conn == nil {
		return errors.New("missing connection, please dial before send data")
	}
	_, err := c.conn.Write(msg)
	return err
}

// OnClose implements TCPClient
func (c *Client) Close() {
	c.conn.Close()
	c.closeSignal <- true
}

func (c *Client) onMessage() {
	for {
		tmp := make([]byte, 256)
		_, err := c.conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				log.Println("read error:", err)
			}
			break
		}
		c.receive <- tmp
	}
}
