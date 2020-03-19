package queue

import (
	"github.com/adjust/rmq"
	"math/big"
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