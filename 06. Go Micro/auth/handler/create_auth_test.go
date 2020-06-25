package handler

import (
	"auth/dao/user"
	"auth/model"
	proto "auth/proto/auth"
	"auth/tool/jwt"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"testing"
	"time"
)

type returns []interface{}
type method string

type CreateAuthTest struct {
	UserId        string
	UserPw        string
	Name          string
	PhoneNumber   string
	Email         string
	Introduction  string
	Authorization string
	ExpectCode    int64
	ExpectMethods map[method]returns
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
	test.Authorization = c.Authorization

	if c.UserId == None 		{ test.UserId = "" } 		else if c.UserId == "" 		  { test.UserId = DefaultUserId }
	if c.UserPw == None 		{ test.UserPw = "" } 		else if c.UserPw == "" 		  { test.UserPw = DefaultUserPw }
	if c.Name == None 		 	{ test.Name = "" } 			else if c.Name == "" 		  { test.Name = DefaultName }
	if c.PhoneNumber == None  	{ test.PhoneNumber = "" }	else if c.PhoneNumber == ""   { test.PhoneNumber = DefaultPN }
	if c.Introduction == None	{ test.Introduction = "" } 	else if c.Introduction == ""  { test.Introduction = "" }
	if c.Email == None 			{ test.Email = "" } 		else if c.Email == "" 		  { test.Email = DefaultEmail }
	if c.Authorization == None  { test.Authorization = "" } else if c.Authorization == "" {
		test.Authorization = jwt.GenerateDuplicateCertJWTNoReturnErr(test.UserId, time.Hour)
	}

	if _, ok := test.ExpectMethods["Insert"]; !ok {
		return
	}

	switch auth := test.ExpectMethods["Insert"][0]; auth.(type) {
	case *model.Auth:
		test.setAuthTupleContext(auth.(*model.Auth), id)
		id++
	case nil:
	}

	return
}

func (c CreateAuthTest) setAuthTupleContext(auth *model.Auth, id uint) {
	auth.ID = id
	auth.UserId = c.UserId
	auth.UserPw = c.UserPw
	auth.Status = user.CreatePending
	auth.CreatedAt = time.Now()
	auth.UpdatedAt = time.Now()
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
	for name, returns := range c.ExpectMethods {
		c.onMethod(string(name), returns)
	}
}

func (c CreateAuthTest) onMethod(method string, returns returns) {
	switch method {
	case "Insert":
		mockStore.On("Insert", &model.Auth{
			UserId: c.UserId,
			UserPw: c.UserPw,
			Status: user.CreatePending,
		}).Return(returns...)
	case "Commit":
		mockStore.On("Commit").Return(returns...)
	case "Rollback":
		mockStore.On("Rollback").Return(returns...)
	default:
		log.Fatalf("%s method cannot be on booked\n", method)
	}
	return
}

func TestAuthCreateManySuccess(t *testing.T) {
	setUpEnv()
	req := &proto.CreateAuthRequest{}
	resp := &proto.CreateAuthResponse{}
	var tests []CreateAuthTest

	forms := []CreateAuthTest {
		{
			UserId:        "testId1",
			ExpectMethods: map[method]returns{
				"Insert": {new(model.Auth), nil},
				"Commit": {new(gorm.DB)},
			},
			ExpectCode:    int64(http.StatusCreated),
		}, {
			UserId:        "testId2",
			ExpectMethods: map[method]returns{
				"Insert": {new(model.Auth), nil},
				"Commit": {new(gorm.DB)},
			},
			ExpectCode:    int64(http.StatusCreated),
		}, {
			UserId:        "testId3",
			ExpectMethods: map[method]returns{
				"Insert": {new(model.Auth), nil},
				"Commit": {&gorm.DB{}},
			},
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

func TestAuthCreateUserIdDuplicateError(t *testing.T) {
	setUpEnv()
	req := &proto.CreateAuthRequest{}
	resp := &proto.CreateAuthResponse{}
	var tests []CreateAuthTest

	var forms = []CreateAuthTest{
		{
			UserId:        "testId1",
			Email:         "jinhong0719@naver.com",
			Authorization: None,
			ExpectCode:    StatusUserIdDuplicate,
		}, {
			UserId:        "testId2",
			Email:         "jinhong0719@naver.com",
			Authorization: jwt.GenerateDuplicateCertJWTNoReturnErr("testId2", "", time.Hour),
			ExpectCode:    StatusUserIdDuplicate,
		}, {
			UserId:        "testId2",
			Email:         "jinhong0719@naver.com",
			Authorization: jwt.GenerateDuplicateCertJWTNoReturnErr("testId1", "", time.Hour),
			ExpectCode:    StatusUserIdDuplicate,
		}, {
			UserId:        "testId2",
			Email:         "jinhong0719@naver.com",
			Authorization: jwt.GenerateDuplicateCertJWTNoReturnErr("", "jinhong0719@naver.com", time.Hour),
			ExpectCode:    StatusUserIdDuplicate,
		}, {
			UserId:        "testId2",
			Email:         "jinhong0719@naver.com",
			Authorization: jwt.GenerateDuplicateCertJWTNoReturnErr("testId2", "jinhong0719@naver.fake", time.Hour),
			ExpectCode:    StatusUserIdDuplicate,
		}, {
			UserId:        "testId2",
			Email:         "jinhong0719@naver.com",
			Authorization: jwt.GenerateDuplicateCertJWTNoReturnErr("testId2", "jinhong0719@naver.com", time.Hour),
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
