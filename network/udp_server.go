package network

import (
	"fmt"
	"log"
	"net"
)

type IBroadcaster interface {
	Broadcast(message []byte) error
	OnBroadcaster()
	Send(addr string, message []byte) error
}

type UDPServer struct {
	port     int
	dialPort int
}

func NewUDPServer(port, dialPort int) IBroadcaster {
	return &UDPServer{
		port:     port,
		dialPort: dialPort,
	}
}

func (b *UDPServer) Broadcast(message []byte) error {
	broadcastAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("255.255.255.255:%d", b.dialPort))
	if err != nil {
		return err
	}

	udpConn, err := net.DialUDP("udp", nil, broadcastAddr)
	if err != nil {
		return err
	}

	_, err = udpConn.Write(message)
	if err != nil {
		return err
	}

	return nil
}

func (b *UDPServer) Send(addr string, message []byte) error {
	broadcastAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	udpConn, err := net.DialUDP("udp", nil, broadcastAddr)
	if err != nil {
		return err
	}

	_, err = udpConn.Write(message)
	if err != nil {
		return err
	}

	return nil
}

func (b *UDPServer) OnBroadcaster() {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", b.port))
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	buf := make([]byte, 1024)

	for {
		i, udpAddr, e := conn.ReadFromUDP(buf)
		if e != nil {
			continue
		}
		fmt.Printf("from %v, reading:%s\n", udpAddr, buf[:i])

		// write back to server
		conn.WriteToUDP([]byte("hello"), udpAddr)
	}
}
