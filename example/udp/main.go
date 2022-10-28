package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/pascallin/go-sak/network"
)

const port = 9099
const dialport = 9098

func d() {
	b := network.NewUDPServer(port, dialport)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go broadcast(b)
	go b.OnBroadcaster()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("Shutting down server...")
}

func broadcast(b network.IBroadcaster) {
	for {
		fmt.Println("try broadcaster")
		err := b.Broadcast([]byte("hello world again"))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("sent hello world")

		time.Sleep(time.Second * 10)
	}
}
