package network

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"

	uuid "github.com/satori/go.uuid"
)

type TCPServer interface {
	Listen(port int) error
	OnRegister() <-chan TCPServerConn
	OnUnregister() <-chan TCPServerConn
	OnReceive() <-chan string
}

type TCPBasicServer struct {
	clients    map[ConnID]TCPServerConn
	register   chan TCPServerConn
	unregister chan TCPServerConn
	receiver   chan string
}

// OnReceive implements TCPServer
func (s *TCPBasicServer) OnReceive() <-chan string {
	return s.receiver
}

// OnUnregister implements TCPServer
func (s *TCPBasicServer) OnUnregister() <-chan TCPServerConn {
	return s.unregister
}

// OnRegister implements TCPServer
func (s *TCPBasicServer) OnRegister() <-chan TCPServerConn {
	return s.register
}

func NewTCPServer() TCPServer {
	return &TCPBasicServer{
		make(map[ConnID]TCPServerConn),
		make(chan TCPServerConn),
		make(chan TCPServerConn),
		make(chan string),
	}
}

func (s *TCPBasicServer) Listen(port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Println("tcp server listener error:", err)
		return err
	}
	go s.listenRegister(listener)
	return nil
}

func (s *TCPBasicServer) listenRegister(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("tcp server accept error", err)
			continue
		}
		go s.registerConnClient(conn)
	}
}

func (s *TCPBasicServer) registerConnClient(conn net.Conn) {
	// create server side client
	id := uuid.NewV4()
	client := TCPServerConn{
		ID:   id,
		conn: conn,
	}
	s.clients[id] = client

	s.register <- client

	defer func() {
		client.close()
		go s.unregisterConnClient(id)
	}()
	for {
		message, err := bufio.NewReader(client.conn).ReadString('\n')
		if err != nil || err == io.EOF {
			continue
		}
		s.receiver <- message
	}
}

func (s *TCPBasicServer) unregisterConnClient(id ConnID) {
	delete(s.clients, id)
	s.unregister <- s.clients[id]
}

func (s *TCPBasicServer) BroadcastMessage(message string) error {
	for _, client := range s.clients {
		err := client.send(message)
		if err != nil {
			return err
		}
	}
	return nil
}
