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

func main() {
	port := 9000
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	server := network.NewTCPServer()
	err := server.Listen(port)
	if err != nil {
		log.Fatal(err)
	}

	var cid network.ConnID
	go func() {
		for {
			select {
			case conn := <-server.OnRegister():
				log.Println("server register ", conn)
				cid = conn.ID
			case conn := <-server.OnUnregister():
				log.Println("server unregister ", conn)
			case msg := <-server.OnReceive():
				log.Println("server received ", string(msg))
				err := server.Send(cid, []byte("hello back"))
				if err != nil {
					log.Println(err)
				}
			}
		}
	}()

	log.Println("listening server...")

	client := network.NewTCPClient()
	client.Dia(fmt.Sprintf("localhost:%d", port))

	go func() {
		for msg := range client.OnReceive() {
			log.Println("client received ", string(msg))
		}
	}()

	go func() {
		for {
			hello := "hello world"
			log.Println("client send", hello)
			err := client.Send([]byte(hello))
			if err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Second * 3)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("Shutting down server...")
}
