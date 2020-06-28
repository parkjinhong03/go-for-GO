package subscriber

import (
	"auth/dao"
	"github.com/stretchr/testify/mock"
)

var mockStore mock.Mock
var authId uint
var h *msgHandler

const (
	none = "none"
	defaultUserId = "TestId"
	defaultUserPw = "TestPw"
	defaultName = "박진홍"
	defaultPN = "01088378347"
	defaultEmail = "jinhong0719@naver.com"
)

func init() {
	mockStore = mock.Mock{}
	adc := dao.NewAuthDAOCreator(nil)
	h = NewMsgHandler(adc)
}

func setUp() {
	mockStore = mock.Mock{}
	authId = 0
}