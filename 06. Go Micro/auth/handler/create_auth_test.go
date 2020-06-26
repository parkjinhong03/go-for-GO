package handler

import (
	"auth/dao/user"
	"auth/model"
	proto "auth/proto/auth"
	"auth/tool/jwt"
	"github.com/google/uuid"
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
	XRequestId    string
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
	test.XRequestId = c.XRequestId

	if c.UserId == None 		{ test.UserId = "" } 		else if c.UserId == "" 		  { test.UserId = DefaultUserId }
	if c.UserPw == None 		{ test.UserPw = "" } 		else if c.UserPw == "" 		  { test.UserPw = DefaultUserPw }
	if c.Name == None 		 	{ test.Name = "" } 			else if c.Name == "" 		  { test.Name = DefaultName }
	if c.PhoneNumber == None  	{ test.PhoneNumber = "" }	else if c.PhoneNumber == ""   { test.PhoneNumber = DefaultPN }
	if c.Introduction == None	{ test.Introduction = "" } 	else if c.Introduction == ""  { test.Introduction = "" }
	if c.Email == None 			{ test.Email = "" } 		else if c.Email == "" 		  { test.Email = DefaultEmail }
	if c.XRequestId == None 	{ test.XRequestId = "" }	else if c.XRequestId == ""	  { test.XRequestId = uuid.New().String() }
	if c.Authorization == None  { test.Authorization = "" } else if c.Authorization == "" {
		test.Authorization = jwt.GenerateDuplicateCertJWTNoReturnErr(test.UserId, test.Email, time.Hour)
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

func (c CreateAuthTest) setRequestContext(req *proto.BeforeCreateAuthRequest) {
	req.UserId = c.UserId
	req.UserPw = c.UserPw
	req.Name = c.Name
	req.Email = c.Email
	req.PhoneNumber = c.PhoneNumber
	req.Introduction = c.Introduction
	req.Authorization = c.Authorization
	req.XRequestID = c.XRequestId
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
	req := &proto.BeforeCreateAuthRequest{}
	resp := &proto.BeforeCreateAuthResponse{}
	var tests []CreateAuthTest

	forms := []CreateAuthTest {
		{
			UserId:        "testId1",
			ExpectMethods: map[method]returns{
				"Insert": {new(model.Auth), nil},
				"Commit": {new(gorm.DB)},
			},
			ExpectCode:    http.StatusCreated,
		}, {
			UserId:        "testId2",
			ExpectMethods: map[method]returns{
				"Insert": {new(model.Auth), nil},
				"Commit": {new(gorm.DB)},
			},
			ExpectCode:    http.StatusCreated,
		}, {
			UserId:        "testId3",
			ExpectMethods: map[method]returns{
				"Insert": {new(model.Auth), nil},
				"Commit": {&gorm.DB{}},
			},
			ExpectCode:    http.StatusCreated,
		},
	}

	for _, form := range forms {
		tests = append(tests, form.createTestFromForm())
	}

	for _, test := range tests {
		test.setRequestContext(req)
		test.onExpectMethods()
		_ = h.BeforeCreateAuth(ctx, req, resp)
		assert.Equal(t, test.ExpectCode, resp.Status)
	}

	mockStore.AssertExpectations(t)
}

func TestBeforeCreateAuthUserIdDuplicateError(t *testing.T) {
	setUpEnv()
	req := &proto.BeforeCreateAuthRequest{}
	resp := &proto.BeforeCreateAuthResponse{}
	var tests []CreateAuthTest

	var forms = []CreateAuthTest{
		{
			UserId:        "testId2",
			Email:         "jinhong0719@naver.com",
			Authorization: jwt.GenerateDuplicateCertJWTNoReturnErr("testId2", "", time.Hour),
			ExpectCode:    StatusEmailDuplicate,
		}, {
			UserId:        "testId2",
			Email:         "jinhong0719@naver.com",
			Authorization: jwt.GenerateDuplicateCertJWTNoReturnErr("testId1", "jinhong0719@naver.com", time.Hour),
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
			ExpectCode:    StatusEmailDuplicate,
		}, {
			UserId:        "testId2",
			Email:         "jinhong0719@naver.com",
			Authorization: jwt.GenerateDuplicateCertJWTNoReturnErr("testId2", "jinhong0719@naver.com", time.Hour),
			ExpectCode:    http.StatusCreated,
			ExpectMethods: map[method]returns{
				"Insert": {new(model.Auth), nil},
				"Commit": {&gorm.DB{}},
			},
		},
	}

	for _, form := range forms {
		tests = append(tests, form.createTestFromForm())
	}

	for _, test := range tests {
		test.setRequestContext(req)
		test.onExpectMethods()
		_ = h.BeforeCreateAuth(ctx, req, resp)
		assert.Equalf(t, test.ExpectCode, resp.Status, "assert error test case: %v\n", test)
	}

	mockStore.AssertExpectations(t)
}

func TestAuthCreateInsertBadRequest(t *testing.T) {
	setUpEnv()
	req := &proto.BeforeCreateAuthRequest{}
	resp := &proto.BeforeCreateAuthResponse{}
	var tests []CreateAuthTest

	forms := []CreateAuthTest{
		{
			UserId:     None,
			ExpectCode: http.StatusBadRequest,
		}, {
			UserPw:     None,
			ExpectCode: http.StatusBadRequest,
		}, {
			UserId:     None,
			UserPw:     None,
			ExpectCode: http.StatusBadRequest,
		}, {
			Authorization: None,
			ExpectCode: http.StatusBadRequest,
		}, {
			XRequestId: None,
			ExpectCode: http.StatusBadRequest,
		}, {
			UserPw:     "qwe",
			ExpectCode: http.StatusBadRequest,
		}, {
			UserId:     "qwe",
			ExpectCode: http.StatusBadRequest,
		}, {
			UserId:     "qwe",
			UserPw:     "qwe",
			ExpectCode: http.StatusBadRequest,
		}, {
			UserId:     "qwerqwerqwerqwerqwer",
			ExpectCode: http.StatusBadRequest,
		}, {
			UserPw:     "qwerqwerqwerqwerqwer",
			ExpectCode: http.StatusBadRequest,
		}, {
			Name:       "박진홍입니다",
			ExpectCode: http.StatusBadRequest,
		}, {
			Name:       "응",
			ExpectCode: http.StatusBadRequest,
		}, {
			PhoneNumber: "0108837834701088378347",
			ExpectCode: http.StatusBadRequest,
		}, {
			Email: "itIsNotEmailFormat",
			ExpectCode: http.StatusBadRequest,
		}, {
			Email: "itIsSoVeryTooLongEmail@naver.com",
			ExpectCode: http.StatusBadRequest,
		}, {
			ExpectCode: http.StatusBadRequest,

		},
	}

	for _, form := range forms {
		tests = append(tests, form.createTestFromForm())
	}

	for _, test := range tests {
		test.setRequestContext(req)
		test.onExpectMethods()
		_ = h.BeforeCreateAuth(ctx, req, resp)
		assert.Equal(t, test.ExpectCode, resp.Status)
	}

	mockStore.AssertExpectations(t)
}
