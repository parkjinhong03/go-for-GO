package handler

import (
	"auth/dao"
	"auth/dao/user"
	"auth/model"
	proto "auth/proto/auth"
	"auth/tool/validator"
	"context"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log"
	"net/http"
	"testing"
)

var mockStore mock.Mock
var ctx context.Context
var h *auth

type Test struct {
	UserId        string
	UserPw        string
	Name          string
	PhoneNumber   string
	Email         string
	Introduction  string
	ExpectCode    int64
	ExpectMethods []string
}

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
	h = NewAuth(adc, validate)
}

func setUpEnv() (req *proto.CreateAuthRequest, rsp *proto.CreateAuthResponse) {
	mockStore = mock.Mock{}
	user.AuthArr = nil
	req = &proto.CreateAuthRequest{}
	rsp = &proto.CreateAuthResponse{}
	return
}

func CreateTestFromForm(form Test) (test Test) {
	test.UserId = form.UserId
	test.UserPw = form.UserPw
	test.Name = form.Name
	test.PhoneNumber = form.PhoneNumber
	test.Introduction = form.Introduction
	test.Email = form.Email
	test.ExpectMethods = form.ExpectMethods
	test.ExpectCode = form.ExpectCode

	if form.UserId == None 		 { test.UserId = "" } 		else if form.UserId == "" 		{ test.UserId = DefaultUserId }
	if form.UserPw == None 		 { test.UserPw = "" } 		else if form.UserPw == "" 		{ test.UserPw = DefaultUserPw }
	if form.Name == None 		 { test.Name = "" } 		else if form.Name == "" 		{ test.Name = DefaultName }
	if form.PhoneNumber == None  { test.PhoneNumber = "" }  else if form.PhoneNumber == ""  { test.PhoneNumber = DefaultPN }
	if form.Introduction == None { test.Introduction = "" } else if form.Introduction == "" { test.Introduction = "" }
	if form.Email == None 		{ test.Email = "" } 		else if form.Email == "" 		{ test.Email = DefaultEmail }

	fmt.Println(test)
	return
}

func setRequestContext(req *proto.CreateAuthRequest, test Test) {
	req.UserId = test.UserId
	req.UserPw = test.UserPw
	req.Name = test.Name
	req.Email = test.Email
	req.PhoneNumber = test.PhoneNumber
	req.Introduction = test.Introduction
}

func onExpectMethods(test Test) {
	for _, expectMethod := range test.ExpectMethods {
		switch expectMethod {
		case "Insert":
			mockStore.On("Insert", &model.Auth{
				UserId: test.UserId,
				UserPw: test.UserPw,
				Status: user.CreatePending,
			}).Return(&model.Auth{}, errors.New(""))
		case "Commit":
			mockStore.On("Commit").Return(&gorm.DB{})
		case "Rollback":
			mockStore.On("Rollback").Return(&gorm.DB{})
		}
	}
}

func TestAuthCreateManySuccess(t *testing.T) {
	req, resp := setUpEnv()
	var tests []Test

	forms := []Test {
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
		tests = append(tests, CreateTestFromForm(form))
	}

	for _, test := range tests {
		setRequestContext(req, test)
		onExpectMethods(test)
		_ = h.CreateAuth(ctx, req, resp)
		assert.Equal(t, test.ExpectCode, resp.Status)
	}

	mockStore.AssertExpectations(t)
}

func TestAuthCreateUserIdAndStudentNumDuplicateError(t *testing.T) {
	req, resp := setUpEnv()
	var tests []Test

	var forms = []Test{
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
		tests = append(tests, CreateTestFromForm(form))
	}

	for _, test := range tests {
		setRequestContext(req, test)
		onExpectMethods(test)
		_ = h.CreateAuth(ctx, req, resp)
		assert.Equal(t, test.ExpectCode, resp.Status)
	}

	mockStore.AssertExpectations(t)
}

func TestAuthCreateInsertBadRequest(t *testing.T) {
	req, resp := setUpEnv()
	var tests []Test

	forms := []Test{
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
			Email: "itIsSoLongEmail@naver.com",
			ExpectCode: int64(http.StatusBadRequest),
		},
	}

	for _, form := range forms {
		tests = append(tests, CreateTestFromForm(form))
	}
	fmt.Println(tests)
	defer func() {
		fmt.Println(req)
	}()

	for _, test := range tests {
		setRequestContext(req, test)
		onExpectMethods(test)
		_ = h.CreateAuth(ctx, req, resp)
		assert.Equal(t, test.ExpectCode, resp.Status)
	}

	mockStore.AssertExpectations(t)
}
