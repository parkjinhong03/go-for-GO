package handler

import (
	"auth/dao"
	"auth/dao/user"
	"auth/model"
	proto "auth/proto/auth"
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

var mockStore mock.Mock
var ctx context.Context
var h *auth

func init() {
	ctx = context.WithValue(context.Background(), "env", "test")
	ctx = context.WithValue(ctx, "mockStore", &mockStore)
	adc := dao.NewAuthDAOCreator(nil)
	validate := validator.New()
	h = NewAuth(adc, validate)
}

func setup() (req *proto.CreateAuthRequest, rsp *proto.CreateAuthResponse) {
	mockStore = mock.Mock{}
	user.AuthArr = nil
	req = &proto.CreateAuthRequest{}
	rsp = &proto.CreateAuthResponse{}
	return
}

func TestAuthCreateManySuccess(t *testing.T) {
	req, resp := setup()

	tests := []struct {
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

	for _, test := range tests {
		req.UserId = test.id
		req.UserPw = test.pw

		mockStore.On("Insert", &model.Auth{
			UserId: test.id,
			UserPw: test.pw,
		}).Return(&model.Auth{}, errors.New(""))
		mockStore.On(test.afterMethod).Return(&gorm.DB{})

		_ = h.CreateAuth(ctx, req, resp)
		assert.Equal(t, test.expectCode, resp.Status)
	}

	mockStore.AssertExpectations(t)
}

func TestAuthCreateUserIdDuplicateError(t *testing.T) {
	req, resp := setup()

	tests := []struct {
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
			id:          "testId1",
			pw:          "testPw1",
			afterMethod: "Rollback",
			expectCode:  int64(StatusUserIdDuplicate),
		}, {
			id:          "testId3",
			pw:          "testPw2",
			afterMethod: "Commit",
			expectCode:  int64(http.StatusCreated),
		},
	}

	for _, test := range tests {
		req.UserId = test.id
		req.UserPw = test.pw

		mockStore.On("Insert", &model.Auth{
			UserId: test.id,
			UserPw: test.pw,
		}).Return(&model.Auth{}, errors.New(""))
		mockStore.On(test.afterMethod).Return(&gorm.DB{})

		_ = h.CreateAuth(ctx, req, resp)
		assert.Equal(t, test.expectCode, resp.Status)
	}

	mockStore.AssertExpectations(t)
}

func TestAuthCreateInsertBadRequest(t *testing.T) {
	req, resp := setup()

	tests := []struct {
		id          string
		pw          string
		expectCode  int64
	}{
		{
			id:          "",
			pw:          "testPw1",
			expectCode:  int64(http.StatusBadRequest),
		}, {
			id:          "testId1",
			pw:          "",
			expectCode:  int64(http.StatusBadRequest),
		}, {
			id:          "",
			pw:          "",
			expectCode:  int64(http.StatusBadRequest),
		}, {
			id:          "qwe",
			pw:          "qewrqewr",
			expectCode:  int64(http.StatusBadRequest),
		}, {
			id:          "qwe",
			pw:          "qwe",
			expectCode:  int64(http.StatusBadRequest),
		}, {
			id:          "qwerqwerqwerqwerqwer",
			pw:          "qweqrwe",
			expectCode:  int64(http.StatusBadRequest),
		}, {
			id:          "qwe",
			pw:          "qwerqwerqwerqwerqwer",
			expectCode:  int64(http.StatusBadRequest),
		},
	}

	for _, test := range tests {
		req.UserId = test.id
		req.UserPw = test.pw
		_ = h.CreateAuth(ctx, req, resp)
		assert.Equal(t, test.expectCode, resp.Status)
	}
	mockStore.AssertExpectations(t)
}
