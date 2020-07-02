package handler

import "github.com/stretchr/testify/mock"

type method string
type returns []interface{}

const (
	none = "none"
	defaultEmail = "jinhong0719@naver.com"
)

var (
	mockStore = mock.Mock{}
)

func init() {
	mockStore = mock.Mock{}
}