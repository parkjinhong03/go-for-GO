package queue

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/adjust/rmq"
	"math/big"
	"strconv"
	"time"
)

// rmq.Queue를 필드로 가지고 있는 메시지 큐 객체를 정의한다.
type RedisQueue struct {
	Queue rmq.Queue
	name string
	callback func(Message) error
}

// 128 비트로 표현할 수 있는 10진수 중 한가지 수를 랜덤으로 얻기 위한 big.Int 객체 선언
var serialNumberLimit = new(big.Int).Lsh(big.NewInt(1), 128)

// Redis Queue 서버에 접속에 새로운 메시지 큐를 생성하여 반환하는 함수이다.
func NewRedisQueue(address, queueName string) *RedisQueue {
	// rmq.OpenConnection 함수를 이용해 매개변수로 받은 주소에 바인딩 되어있는 Redis Queue 서버에 tcp로 연결한다.
	conn := rmq.OpenConnection("test service", "tcp", address, 1)
	rq := conn.OpenQueue(queueName)

	return &RedisQueue{
		Queue:    rq,
		name:     queueName,
		callback: nil,
	}
}

func (r *RedisQueue) Add(messageName string, payload []byte) error {
	m := Message{Name: messageName, Payload: string(payload)}
	return r.AddMessage(m)
}

func (r *RedisQueue) AddMessage(message Message) error {
	// 전역 변수로 선언한 serialNumberLimit를 이용해서 128비트 짜리의 랜덤 정수 값을 생성한다.
	serialNumber, _ := rand.Int(rand.Reader, serialNumberLimit)
	// 나노초 단위의 현재 시간을 덧붙여서 serialNumber를 고유값으로 만든 다음 message.ID에 대입한다.
	message.ID = strconv.Itoa(time.Now().Nanosecond()) + serialNumber.String()

	messageByte, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("cannot marshal message to byte: %v", err)
	}

	fmt.Printf("add event to queue: %s", string(messageByte))
	// json.Marshal 함수로 인코딩한 message를 Queue.PublishBytes 메서드에 전달하여 큐에 게시한다.
	if !r.Queue.PublishBytes(messageByte) {
		return fmt.Errorf("cannot add message to queue: %v", err)
	}

	return nil
}