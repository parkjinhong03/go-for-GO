package main

import (
	"building-microservices-with-go/eda/queue"
	"encoding/json"
	"log"
	"runtime"
	"time"
)

type Request struct {
	Name string `json:"name"`
}

func main() {
	log.Println("Starting Worker")

	q := queue.NewRedisQueue("localhost:6379", "test_queue", "reader")

	q.StartConsuming(5, 10*time.Second, func(message queue.Message) error {
		var request Request
		if err := json.Unmarshal([]byte(message.Payload), &request); err != nil {
			return err
		}
		return nil
	})

	runtime.Goexit()
}