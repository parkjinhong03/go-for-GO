package broker

import (
	"auth/subscriber"
	topic "auth/topic/golang"
	br "github.com/micro/go-micro/v2/broker"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/server"
	"github.com/micro/go-plugins/broker/rabbitmq/v2"
)

func RabbitMQInitializer(s server.Server, as *subscriber.Auth) func() error {
	return func() (err error) {
		brk := s.Options().Broker

		if err = brk.Connect(); err != nil { return }
		_, err = brk.Subscribe(topic.CreateAuthEventTopic, as.CreateAuth,
			br.Queue(topic.CreateAuthEventTopic), // Queue 정적 이름 설정
			br.DisableAutoAck(), // Ack를 수동으로 실행하게 설정
			rabbitmq.DurableQueue()) // Queue 연결을 종료해도 삭제X 설정
		if err != nil { return }

		_, err = brk.Subscribe(topic.ChangeAuthStatusEventTopic, as.ChangeAuthStatus,
			br.Queue(topic.ChangeAuthStatusEventTopic), // Queue 정적 이름 설정
			br.DisableAutoAck(), // Ack를 수동으로 실행하게 설정
			rabbitmq.DurableQueue()) // Queue 연결을 종료해도 삭제X 설정
		if err != nil { return }

		log.Infof("succeed in connecting to broker!! (name: %s | addr: %s)\n",  brk.String(), brk.Address())
		return
	}
}