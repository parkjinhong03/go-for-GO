package handler

import (
	"auth/dao"
	"auth/dao/user"
	"auth/model"
	proto "auth/proto/auth"
	"context"
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

var mockStore mock.Mock
var ctx context.Context
var h auth

func init() {
	ctx = context.WithValue(context.Background(), "env", "test")
	ctx = context.WithValue(ctx, "mockStore", &mockStore)
	h = auth{
		adc: dao.NewAuthDAOCreator(nil),
	}
}

func setup() (req *proto.CreateAuthRequest, rsp *proto.CreateAuthResponse) {
	mockStore = mock.Mock{}
	req = &proto.CreateAuthRequest{}
	rsp = &proto.CreateAuthResponse{}
	return
}

func prepareInsertRequest(req *proto.CreateAuthRequest, id, pw string) {
	req.UserId = id
	req.UserPw = pw

	mockStore.On("Insert", &model.Auth{
		UserId: id,
		UserPw: pw,
	}).Return(&model.Auth{}, errors.New(""))
}

func TestAuthCreateInsertOne(t *testing.T) {
	user.AuthArr = nil
	req, resp := setup()
	prepareInsertRequest(req, "testId1", "testPw1")
	mockStore.On("Commit").Return(&gorm.DB{})

	if err := h.CreateAuth(ctx, req, resp); err != nil {
		resp.Status = http.StatusInternalServerError
	}

	mockStore.AssertExpectations(t)
	assert.Equal(t, int64(http.StatusCreated), resp.Status)
}

func TestAuthCreateInsertTwoAndMore(t *testing.T) {
	user.AuthArr = nil
	req, resp := setup()

	requests := []struct {
		id          string
		pw          string
		afterMethod string
		expectCode  int64
	}{
		{
			id:          "testId1",
			pw:          "testPw1",
			afterMethod: "Commit",
			expectCode:  int64(http.StatusCreated),
		}, {
			id:          "testId2",
			pw:          "testPw1",
			afterMethod: "Commit",
			expectCode:  int64(http.StatusCreated),
		}, {
			id:          "testId3",
			pw:          "testPw2",
			afterMethod: "Commit",
			expectCode:  int64(http.StatusCreated),
		},
	}

	for _, request := range requests {
		prepareInsertRequest(req, request.id, request.pw)
		mockStore.On(request.afterMethod).Return(&gorm.DB{})
		_ = h.CreateAuth(ctx, req, resp)
		assert.Equal(t, request.expectCode, resp.Status)
	}
	mockStore.AssertExpectations(t)
}