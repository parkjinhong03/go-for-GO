package main

import (
	"building-microservices-with-go/eda/queue"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync/atomic"
)

type QueueWriterHandler struct {
	queue *queue.RedisQueue
	queueNum int32
}

func (qw *QueueWriterHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	data, _ := ioutil.ReadAll(r.Body)
	atomic.AddInt32(&qw.queueNum, 1)
	err := qw.queue.Add(fmt.Sprintf("test.product.%d", qw.queueNum), data)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func main() {
	q := queue.NewRedisQueue("localhost:6379", "test_queue", "writer")
	handler := &QueueWriterHandler{
		queue: q,
	}

	http.Handle("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}