// 이 전까지는 어플리케이션들이 동기적으로 통신할 때, 더 안정성있고 높은 성능으로 통신할 수 있는 몇 가지 패턴들을 살펴보았다.
// 하지만 이번 9장 이벤트 주도 아키텍처에서는 여러 어플리케이션들이 비동기 처리로 메시지를 주고 받는 방법에 대해 알아볼 것 이다.
// 비동기식 처리는 종종 풀(pull/queue)과 푸시(push)라는 두 가지 형식의 패턴이 있다.
// 먼저 풀 패턴을 구현해볼 것인데, 풀/큐 메시징은 작업자 프로세스가 실행 중인 경우에 적합한 설계이다.

package queue

import "time"

// Message는 큐에 저장할 메세지의 구조체이다.
type Message struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Payload string `json:"payload"`
}

// Queue는 특정 객체를 큐로 사용하기 위해서 구현해야하는 인터페이스이다.
type Queue interface {
	// 새 매세지를 발행하는 편리한 방법으로 ID는 내부에서 고유한 계산된 값을 생성하여 추가한다.
	Add(messageName string, payload []byte) error
	// name, payload를 Message 구조체에 담아서 전달하는 점을 제외하면 Add와 동일한 메서드이다.
	AddMessage(message Message) error
	// 구독자가 큐에서 메세지를 가져올 수 있도록 하는 메서드이다.
	StartConsuming(size int, pollInterval time.Duration, callback func(Message) error)
}