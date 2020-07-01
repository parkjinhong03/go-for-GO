package main

import (
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"user/adapter/broker"
	"user/adapter/db"
	"user/dao"
	"user/handler"
	"user/subscriber"
	"user/tool/validator"
)

func main() {
	conn, err := db.ConnMysql()
	if err != nil { log.Fatal(err) }
	udc := dao.NewUserDAOCreator(conn)
	validate, err := validator.New()
	if err != nil { log.Fatal(err) }
	rbMQ := broker.ConnRabbitMQ()

	h := handler.NewUser(rbMQ, validate, udc)
	s := subscriber.NewUser()

	service := micro.NewService(
		micro.Name("examples.blog.service.user"),
		micro.Version("latest"),
		micro.Broker(rbMQ),
	)

	brkHandler := func() error {
		brk := service.Options().Broker
		if err := brk.Connect(); err != nil { log.Fatal(err) }
		if _, err := brk.Subscribe(subscriber.CreateUserEventTopic); err != nil { log.Fatal(err) }
		log.Infof("succeed in connecting to broker!! (name: %s | addr: %s)\n",  brk.String(), brk.Address())
		return nil
	}

	service.Init(micro.AfterStart(brkHandler))

	if err = micro.RegisterHandler(service.Server(), h); err != nil {
		log.Fatal(err)
	}
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
