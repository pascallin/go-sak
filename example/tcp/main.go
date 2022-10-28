package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/pascallin/go-sak/network"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	server := network.NewTCPServer()
	server.Listen(8080)

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("Shutting down server...")
}
