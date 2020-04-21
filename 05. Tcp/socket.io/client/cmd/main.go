package main

import (
	"fmt"
	"socket.io/test/client"
	"time"
)

func main() {
	c1 := client.NewSocketClient()
	c2 := client.NewSocketClient()

	_ = c1.Emit("NAME", "qwe")
	_ = c2.Emit("NAME", "asd")

	go func() {
		for {
			select {
			case c := <-c1.Incoming():
				fmt.Println(c.Name, c.Message)
			case c := <-c2.Incoming():
				fmt.Println(c.Name, c.Message)
			}
		}
	}()

	_ = c1.Emit("SEND", "hello")
	_ = c2.Emit("SEND", "hello")

	time.Sleep(time.Second)
}