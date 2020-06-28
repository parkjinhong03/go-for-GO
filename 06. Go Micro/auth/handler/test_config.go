package handler

import (
	"auth/adapter/broker"
	"auth/dao"
	"auth/tool/validator"
	"context"
	"github.com/stretchr/testify/mock"
	"log"
)

var mockStore mock.Mock
var ctx context.Context
var h *auth

const (
	None = "none"
	DefaultUserId = "testId"
	DefaultUserPw = "testPw"
	DefaultName = "박진홍"
	DefaultPN = "01088378347"
	DefaultEmail = "jinhong0719@naver.com"
)

func init() {
	ctx = context.WithValue(context.Background(), "env", "test")
	ctx = context.WithValue(ctx, "mockStore", &mockStore)
	adc := dao.NewAuthDAOCreator(nil)
	validate, err := validator.New()
	if err != nil { log.Fatal(err) }
	b := broker.NewRabbitMQForTest(&mockStore)
	h = NewAuth(b, adc, validate)
}

func setUpEnv() () {
	mockStore = mock.Mock{}
}
