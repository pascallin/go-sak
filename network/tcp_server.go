package network

import (
	"errors"
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
	OnReceive() <-chan []byte
	Send(cid ConnID, msg []byte) error
	Broadcast(msg []byte) error
}

type TCPBasicServer struct {
	connections map[ConnID]TCPServerConn
	register    chan TCPServerConn
	unregister  chan TCPServerConn
	receiver    chan []byte
}

// Broadcast implements TCPServer
func (s *TCPBasicServer) Broadcast(msg []byte) error {
	for _, client := range s.connections {
		err := client.Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

// Send implements TCPServer
func (s *TCPBasicServer) Send(cid uuid.UUID, msg []byte) error {
	_, ok := s.connections[cid]
	if !ok {
		return errors.New("could not found connection id " + cid.String())
	}
	return s.connections[cid].Send(msg)
}

// OnReceive implements TCPServer
func (s *TCPBasicServer) OnReceive() <-chan []byte {
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
		make(chan []byte),
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
	s.connections[id] = client

	s.register <- client

	defer func() {
		client.close()
		go s.unregisterConnClient(id)
	}()
	for {
		tmp := make([]byte, 256)
		_, err := conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				log.Println("read error:", err)
			}
			break
		}
		s.receiver <- tmp
	}
}

func (s *TCPBasicServer) unregisterConnClient(id ConnID) {
	delete(s.connections, id)
	s.unregister <- s.connections[id]
}
