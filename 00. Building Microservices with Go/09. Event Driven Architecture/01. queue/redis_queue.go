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

	fmt.Printf("Add event to queue: %s", string(messageByte))
	// json.Marshal 함수로 인코딩한 message를 Queue.PublishBytes 메서드에 전달하여 큐에 게시한다.
	if !r.Queue.PublishBytes(messageByte) {
		return fmt.Errorf("cannot add message to queue: %v", err)
	}

	return nil
}

// StartConsuming 함수는 매개변수로 받은 콜백 함수를 큐 인스턴스에 등록하고, 메세지 풀링을 시작시킴과 동시에 구독자를 등록시키는 함수이다.
func (r *RedisQueue) StartConsuming(prefetchSize int, pullDuration time.Duration, callback func(Message) error) {
	r.callback = callback
	// Queue.StartConsuming 메서드를 이용하여 이제부터 큐에 쌓인 메세지를 처리하겠다는 신호를 보낸다.
	r.Queue.StartConsuming(prefetchSize, pullDuration)
	// Queue.AddConsumer 메서드를 이용하여 메세지를 처리할 구독자를 등록시킬 수 있다.
	// 이제부터 큐에서 새 메세지가 감지되면 rmq.Delivery 인터페이스를 구현한 두 번째 매개변수의 Consume(delivery rmq.Delivery) 메서드를 실행시킨다.
	r.Queue.AddConsumer("RedisQueue_"+r.name, r)
}

// 위에서 말했다 시피, 새 메세지가 감지되면 호출되는 함수이다.
func (r *RedisQueue) Consume(delivery rmq.Delivery) {
	// delivery.Payload 메서드를 이용하여 메세지의 내용(payload)를 얻을 수 있다.
	fmt.Printf("Got event from redis queue: %s", delivery.Payload())
	var message Message
	if err := json.Unmarshal([]byte(delivery.Payload()), &message); err != nil {
		fmt.Println("Error consuming event, unable to unmarshal event")
		// payload를 message 구조체로 언마샬링 하는데 실패한다면, 메세지를 나중에 처리할 수 있도록 남겨두기 위해 delivery.Reject 메서드를 호출한다.
		delivery.Reject()
		return
	}

	// 그리고 StartConsuming 함수에서 등록했던 callback 함수를 실행시킨다.
	// 이 함수에서는 실질적으로 메세지를 처리하는 비즈니스 로직이 포함되어있다.
	if err := r.callback(message); err != nil {
		fmt.Println("Error consuming event, in processing callback function")
		// 콜백 함수도 실패하였다면 이 메세지 또한 남겨두기 위해 delivery.Reject 메서드를 실행시킨다.
		delivery.Reject()
		return
	}

	// 성공적으로 메세지 처리가 완료되었다면 delivery.Ack 메서드를 실행시켜 큐에서 해당 메세지를 제거시킨다.
	delivery.Ack()
}