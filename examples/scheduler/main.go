package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pascallin/go-sak/pkg/scheduler"
)

type Stu struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	s := scheduler.New(func(s *scheduler.Scheduler, e *scheduler.Event) {
		fmt.Println("Scheduler function with event", e)
	})
	//now time
	now := time.Now().Add(time.Second * 10).Format(time.RFC3339)
	fmt.Println("Nowtime:", now)
	stu := Stu{
		Name: "pascal",
		Age:  18,
	}
	jsonStu, err := json.Marshal(stu)
	if err != nil {
		fmt.Println("Parse json error")
	}
	attachments := []scheduler.Attachment{{
		Name:        "foo",
		ContentType: "json",
		Body:        []byte(jsonStu),
	}}
	s.Schedule(scheduler.NewEvent(now, attachments))
	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
	fmt.Println("exit!")
}
