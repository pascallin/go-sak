package network

import (
	"os"
	"testing"
)

func TestClient_DiaAndSend(t *testing.T) {
	client := NewTCPClient(os.Stdout)
	client.Dia("localhost:8080")

	client.Send("Hello")
}
