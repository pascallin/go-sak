package network

import "testing"

func TestNewServer(t *testing.T) {
	server := NewServer()
	server.Listen(8080)
}
