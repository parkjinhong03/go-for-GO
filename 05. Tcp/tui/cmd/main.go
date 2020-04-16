package main

import (
	"chat.server.com/client"
	"chat.server.com/tui"
	"log"
)

func main() {
	c := client.NewTcpChatClient()
	if err := c.Dial(":8080"); err != nil { log.Fatal(err) }
	go c.Start()
	defer c.Disconnect()

	tui.StartUi(c)
}