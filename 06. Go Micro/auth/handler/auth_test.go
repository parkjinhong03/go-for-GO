package handler

import (
	"auth/dao"
	"auth/model"
	proto "auth/proto/auth"
	"context"
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log"
	"net/http"
	"testing"
)

var mockStore = mock.Mock{}

func setup() (h auth, ctx context.Context, req *proto.CreateAuthRequest, rsp *proto.CreateAuthResponse) {
	adc := dao.NewAuthDAOCreator(nil)
	h = auth{
		adc: adc,
	}
	ctx = context.WithValue(context.Background(), "env", "test")
	ctx = context.WithValue(ctx, "mockStore", &mockStore)
	req = &proto.CreateAuthRequest{}
	rsp = &proto.CreateAuthResponse{}
	return
}

func TestAuthCreateInsertOne(t *testing.T) {
	h, ctx, req, rsp := setup()
	req.UserId = "testId"
	req.UserPw = "testPwd"

	mockStore.On("Insert", &model.Auth{
		UserId: "testId",
		UserPw: "testPwd",
	}).Return(&model.Auth{}, errors.New(""))
	mockStore.On("Commit").Return(&gorm.DB{})

	if err := h.CreateAuth(ctx, req, rsp); err != nil {
		log.Fatal(err)
	}

	mockStore.AssertExpectations(t)
	assert.Equal(t, int64(http.StatusOK), rsp.Status)
}