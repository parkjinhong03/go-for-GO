package main

import (
	"chat.server.com/client"
	"fmt"
	"log"
	"time"
)

func main() {
	c1 := client.NewTcpChatClient()
	c2 := client.NewTcpChatClient()
	if err := c1.Dial("ec2-18-218-105-177.us-east-2.compute.amazonaws.com:8080"); err != nil {
		log.Fatal(err)
	}
	if err := c2.Dial("ec2-18-218-105-177.us-east-2.compute.amazonaws.com:8080"); err != nil {
		log.Fatal(err)
	}

	go c1.Start()
	go c2.Start()

	go func() {
		for {
			select {
			case cmd := <- c1.Incoming():
				fmt.Println(1, cmd.Name, cmd.Message)
			case cmd := <- c2.Incoming():
				fmt.Println(2, cmd.Name, cmd.Message)
			}
		}
	}()

	c1.SetName("qwe")
	c2.SetName("asd")
	c1.SendMessage("hello")

	time.Sleep(5 * time.Second)
	c1.Disconnect()
	c2.Disconnect()
}