package handler

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/uber/jaeger-client-go"
	"log"
	"user/dao"
	"user/tool/validator"
)

type method string
type returns []interface{}

const (
	none = "none"
	defaultEmail = "jinhong0719@naver.com"
)

var (
	mockStore = mock.Mock{}
	h *user
	ctx context.Context
)

func init() {
	ctx = context.WithValue(context.Background(), "env", "test")
	ctx = context.WithValue(ctx, "mockStore", &mockStore)
	mockStore = mock.Mock{}
	validate, err := validator.New()
	if err != nil { log.Fatal(err) }
	udc := dao.NewUserDAOCreator(nil)
	// mq mock 객체 대입 필요
	// tracer 모의 객체 대입 필요
	h = NewUser(nil, validate, udc, &jaeger.Tracer{})
}

func setUpEnv() {
	ctx = context.WithValue(context.Background(), "env", "test")
	ctx = context.WithValue(ctx, "mockStore", &mockStore)
	mockStore = mock.Mock{}
}