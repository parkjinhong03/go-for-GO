package handler

import (
	"auth/dao"
	"auth/dao/user"
	"auth/model"
	proto "auth/proto/auth"
	"auth/tool/validator"
	"context"
	"errors"
	"github.com/bmizerany/assert"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
	"log"
	"net/http"
	"testing"
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
	h = NewAuth(nil, adc, validate)
}

func setUpEnv() () {
	mockStore = mock.Mock{}
	user.AuthArr = nil
}

type CreateAuthTest struct {
	UserId        string
	UserPw        string
	Name          string
	PhoneNumber   string
	Email         string
	Introduction  string
	ExpectCode    int64
	ExpectMethods []string
}

func (c CreateAuthTest) createTestFromForm() (test CreateAuthTest) {
	test.UserId = c.UserId
	test.UserPw = c.UserPw
	test.Name = c.Name
	test.PhoneNumber = c.PhoneNumber
	test.Introduction = c.Introduction
	test.Email = c.Email
	test.ExpectMethods = c.ExpectMethods
	test.ExpectCode = c.ExpectCode

	if c.UserId == None 		{ test.UserId = "" } 		else if c.UserId == "" 		  { test.UserId = DefaultUserId }
	if c.UserPw == None 		{ test.UserPw = "" } 		else if c.UserPw == "" 		  { test.UserPw = DefaultUserPw }
	if c.Name == None 		 	{ test.Name = "" } 			else if c.Name == "" 		  { test.Name = DefaultName }
	if c.PhoneNumber == None  	{ test.PhoneNumber = "" }	else if c.PhoneNumber == ""   { test.PhoneNumber = DefaultPN }
	if c.Introduction == None	{ test.Introduction = "" } 	else if c.Introduction == ""  { test.Introduction = "" }
	if c.Email == None 			{ test.Email = "" } 		else if c.Email == "" 		  { test.Email = DefaultEmail }

	return
}

func (c CreateAuthTest) setRequestContext(req *proto.CreateAuthRequest) {
	req.UserId = c.UserId
	req.UserPw = c.UserPw
	req.Name = c.Name
	req.Email = c.Email
	req.PhoneNumber = c.PhoneNumber
	req.Introduction = c.Introduction
}

func (c CreateAuthTest) onExpectMethods() {
	for _, expectMethod := range c.ExpectMethods {
		c.onMethod(expectMethod)
	}
}

func (c CreateAuthTest) onMethod(method string) {
	switch method {
	case "Insert":
		mockStore.On("Insert", &model.Auth{
			UserId: c.UserId,
			UserPw: c.UserPw,
			Status: user.CreatePending,
		}).Return(&model.Auth{}, errors.New(""))
	case "Commit":
		mockStore.On("Commit").Return(&gorm.DB{})
	case "Rollback":
		mockStore.On("Rollback").Return(&gorm.DB{})
	}
}

func TestAuthCreateManySuccess(t *testing.T) {
	setUpEnv()
	req := &proto.CreateAuthRequest{}
	resp := &proto.CreateAuthResponse{}
	var tests []CreateAuthTest

	forms := []CreateAuthTest {
		{
			UserId:        "testId1",
			ExpectMethods: []string{"Insert", "Commit"},
			ExpectCode:    int64(http.StatusCreated),
		}, {
			UserId:        "testId2",
			ExpectMethods: []string{"Insert", "Commit"},
			ExpectCode:    int64(http.StatusCreated),
		}, {
			UserId:        "testId3",
			ExpectMethods: []string{"Insert", "Commit"},
			ExpectCode:    int64(http.StatusCreated),
		},
	}

	for _, form := range forms {
		tests = append(tests, form.createTestFromForm())
	}

	for _, test := range tests {
		test.setRequestContext(req)
		test.onExpectMethods()
		_ = h.CreateAuth(ctx, req, resp)
		assert.Equal(t, test.ExpectCode, resp.Status)
	}

	mockStore.AssertExpectations(t)
}

func TestAuthCreateUserIdAndStudentNumDuplicateError(t *testing.T) {
	setUpEnv()
	req := &proto.CreateAuthRequest{}
	resp := &proto.CreateAuthResponse{}
	var tests []CreateAuthTest

	var forms = []CreateAuthTest{
		{
			UserId:        "testId1",
			ExpectMethods: []string{"Insert", "Commit"},
			ExpectCode:    int64(http.StatusCreated),
		}, {
			UserId:        "testId1",
			ExpectMethods: []string{"Insert", "Rollback"},
			ExpectCode:    StatusUserIdDuplicate,
		}, {
			UserId:        "testId2",
			ExpectMethods: []string{"Insert", "Commit"},
			ExpectCode:    int64(http.StatusCreated),
		},
	}

	for _, form := range forms {
		tests = append(tests, form.createTestFromForm())
	}

	for _, test := range tests {
		test.setRequestContext(req)
		test.onExpectMethods()
		_ = h.CreateAuth(ctx, req, resp)
		assert.Equal(t, test.ExpectCode, resp.Status)
	}

	mockStore.AssertExpectations(t)
}

func TestAuthCreateInsertBadRequest(t *testing.T) {
	setUpEnv()
	req := &proto.CreateAuthRequest{}
	resp := &proto.CreateAuthResponse{}
	var tests []CreateAuthTest

	forms := []CreateAuthTest{
		{
			UserId:     None,
			ExpectCode: int64(http.StatusBadRequest),
		}, {
			UserPw:     None,
			ExpectCode: int64(http.StatusBadRequest),
		}, {
			UserId:     None,
			UserPw:     None,
			ExpectCode: int64(http.StatusBadRequest),
		}, {
			UserPw:     "qwe",
			ExpectCode: int64(http.StatusBadRequest),
		}, {
			UserId:     "qwe",
			ExpectCode: int64(http.StatusBadRequest),
		}, {
			UserId:     "qwe",
			UserPw:     "qwe",
			ExpectCode: int64(http.StatusBadRequest),
		}, {
			UserId:     "qwerqwerqwerqwerqwer",
			ExpectCode: int64(http.StatusBadRequest),
		}, {
			UserPw:     "qwerqwerqwerqwerqwer",
			ExpectCode: int64(http.StatusBadRequest),
		}, {
			Name:       "박진홍입니다",
			ExpectCode: int64(http.StatusBadRequest),
		}, {
			Name:       "응",
			ExpectCode: int64(http.StatusBadRequest),
		}, {
			PhoneNumber: "0108837834701088378347",
			ExpectCode: int64(http.StatusBadRequest),
		}, {
			Email: "itIsNotEmailFormat",
			ExpectCode: int64(http.StatusBadRequest),
		}, {
			Email: "itIsSoVeryTooLongEmail@naver.com",
			ExpectCode: int64(http.StatusBadRequest),
		},
	}

	for _, form := range forms {
		tests = append(tests, form.createTestFromForm())
	}

	for _, test := range tests {
		test.setRequestContext(req)
		test.onExpectMethods()
		_ = h.CreateAuth(ctx, req, resp)
		assert.Equal(t, test.ExpectCode, resp.Status)
	}

	mockStore.AssertExpectations(t)
}
