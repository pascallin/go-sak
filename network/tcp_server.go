package network

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"

	uuid "github.com/satori/go.uuid"
)

type TCPServer struct {
	clients    map[ConnID]TCPServerConn
	Register   chan TCPServerConn
	Unregister chan TCPServerConn
	Receiver   chan string
}

func NewServer() *TCPServer {
	server := &TCPServer{
		make(map[ConnID]TCPServerConn),
		make(chan TCPServerConn),
		make(chan TCPServerConn),
		make(chan string),
	}

	return server
}

func (s *TCPServer) Listen(port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal("tcp server listener error:", err)
	}
	go s.listenRegister(listener)
}

func (s *TCPServer) listenRegister(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("tcp server accept error", err)
			continue
		}
		go s.registerConnClient(conn)
	}
}

func (s *TCPServer) registerConnClient(conn net.Conn) {
	// create server side client
	id := uuid.NewV4()
	client := TCPServerConn{
		ID:   id,
		conn: conn,
	}
	s.clients[id] = client

	// user register notice
	s.Register <- client

	defer func() {
		client.close()
		go s.unregisterConnClient(id)
	}()
	for {
		message, err := bufio.NewReader(client.conn).ReadString('\n')
		if err != nil || err == io.EOF {
			continue
		}
		s.Receiver <- message
	}
}

func (s *TCPServer) unregisterConnClient(id ConnID) {
	delete(s.clients, id)
	s.Unregister <- s.clients[id]
}

func (s *TCPServer) BroadcastMessage(message string) error {
	for _, client := range s.clients {
		err := client.send(message)
		if err != nil {
			return err
		}
	}
	return nil
}
