package handler

import (
	proto "auth/proto/auth"
	"auth/subscriber"
	"auth/tool/jwt"
	"auth/tool/random"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/micro/go-micro/v2/broker"
	"github.com/stretchr/testify/assert"
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
	ExpectMessage string
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
	test.ExpectMessage = c.ExpectMessage
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
	return
}

//func (c CreateAuthTest) setAuthTupleContext(auth *model.Auth, id uint) {
//	auth.ID = id
//	auth.UserId = c.UserId
//	auth.UserPw = c.UserPw
//	auth.Status = user.CreatePending
//	auth.CreatedAt = time.Now()
//	auth.UpdatedAt = time.Now()
//	id++
//}

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
		c.onMethod(name, returns)
	}
}

func (c CreateAuthTest) onMethod(method method, returns returns) {
	switch method {
	case "Publish":
		header := make(map[string]string)
		header["XRequestId"] = c.XRequestId
		header["MessageId"] = ctx.Value("MessageId").(string)
		mockStore.On("Publish", subscriber.CreateAuthEventTopic, &broker.Message{
			Header: header,
			Body:   nil,
		}).Return(returns...)
	default:
		panic(fmt.Sprintf("%s method cannot be on booked\n", method))
	}
	//switch method {
	//case "Insert":
	//	mockStore.On("Insert", &model.Auth{
	//		UserId: c.UserId,
	//		UserPw: c.UserPw,
	//		Status: user.CreatePending,
	//	}).Return(returns...)
	//case "Commit":
	//	mockStore.On("Commit").Return(returns...)
	//case "Rollback":
	//	mockStore.On("Rollback").Return(returns...)
	//default:
	//	log.Fatalf("%s method cannot be on booked\n", method)
	//}
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
			Email:         "jinhong0719@naver.com",
			Authorization: jwt.GenerateDuplicateCertJWTNoReturnErr("testId1", "jinhong0719@naver.com", time.Hour),
			ExpectMethods: map[method]returns{
				"Publish": {nil},
			},
			ExpectCode:    http.StatusCreated,
			ExpectMessage: MessageAuthCreated,
		}, {
			UserId:        "testId1",
			Email:         "richimous0719@naver.com",
			Authorization: jwt.GenerateDuplicateCertJWTNoReturnErr("testId1", "richimous0719@naver.com", time.Hour),
			ExpectMethods: map[method]returns{
				"Publish": {nil},
			},
			ExpectCode:    http.StatusCreated,
			ExpectMessage: MessageAuthCreated,
		},
	}

	for _, form := range forms {
		tests = append(tests, form.createTestFromForm())
	}

	for _, test := range tests {
		ctx = context.WithValue(ctx, "MessageId", random.GenerateString(32))
		test.setRequestContext(req)
		test.onExpectMethods()
		_ = h.BeforeCreateAuth(ctx, req, resp)
		assert.Equalf(t, test.ExpectCode, resp.Status, "status assert error test case: %v\n", test)
		assert.Equalf(t, test.ExpectMessage, resp.Message, "message assert error test case: %v\n", test)
	}

	mockStore.AssertExpectations(t)
}

func TestBeforeCreateAuthUserIdDuplicateError(t *testing.T) {
	setUpEnv()
	req := &proto.BeforeCreateAuthRequest{}
	resp := &proto.BeforeCreateAuthResponse{}
	var tests []CreateAuthTest

	var forms = []CreateAuthTest{{
			UserId:        "testId2",
			Email:         "jinhong0719@naver.com",
			Authorization: jwt.GenerateDuplicateCertJWTNoReturnErr("testId2", "", time.Hour),
			ExpectCode:    StatusEmailDuplicate,
			ExpectMessage: MessageEmailDuplicate,
		}, {
			UserId:        "testId2",
			Email:         "jinhong0719@naver.com",
			Authorization: jwt.GenerateDuplicateCertJWTNoReturnErr("testId1", "jinhong0719@naver.com", time.Hour),
			ExpectCode:    StatusUserIdDuplicate,
			ExpectMessage: MessageUserIdDuplicate,
		}, {
			UserId:        "testId2",
			Email:         "jinhong0719@naver.com",
			Authorization: jwt.GenerateDuplicateCertJWTNoReturnErr("", "jinhong0719@naver.com", time.Hour),
			ExpectCode:    StatusUserIdDuplicate,
			ExpectMessage: MessageUserIdDuplicate,
		}, {
			UserId:        "testId2",
			Email:         "jinhong0719@naver.com",
			Authorization: jwt.GenerateDuplicateCertJWTNoReturnErr("testId2", "jinhong0719@naver.fake", time.Hour),
			ExpectCode:    StatusEmailDuplicate,
			ExpectMessage: MessageEmailDuplicate,
		},
	}

	for _, form := range forms {
		tests = append(tests, form.createTestFromForm())
	}

	for _, test := range tests {
		test.setRequestContext(req)
		test.onExpectMethods()
		_ = h.BeforeCreateAuth(ctx, req, resp)
		assert.Equalf(t, test.ExpectCode, resp.Status, "status assert error test case: %v\n", test)
		assert.Equalf(t, test.ExpectMessage, resp.Message, "message assert error test case: %v\n", test)
	}

	mockStore.AssertExpectations(t)
}

func TestBeforeCreateAuthForbidden(t *testing.T) {
	setUpEnv()
	req := &proto.BeforeCreateAuthRequest{}
	resp := &proto.BeforeCreateAuthResponse{}
	var tests []CreateAuthTest

	var forms = []CreateAuthTest{
		{
			UserId:        "testId2",
			Email:         "jinhong0719@naver.com",
			Authorization: "ThisIsInvalidAuthorizationString",
			ExpectCode:    http.StatusForbidden,
		},
	}

	for _, form := range forms {
		tests = append(tests, form.createTestFromForm())
	}

	for _, test := range tests {
		test.setRequestContext(req)
		test.onExpectMethods()
		_ = h.BeforeCreateAuth(ctx, req, resp)
		assert.Equalf(t, test.ExpectCode, resp.Status, "status assert error test case: %v\n", test)
	}

	mockStore.AssertExpectations(t)
}

func TestBeforeAuthCreateInsertBadRequest(t *testing.T) {
	setUpEnv()
	req := &proto.BeforeCreateAuthRequest{}
	resp := &proto.BeforeCreateAuthResponse{}
	var tests []CreateAuthTest

	forms := []CreateAuthTest{
		{
			UserId: None,
		}, {
			UserPw: None,
		}, {
			UserId: None,
			UserPw: None,
		}, {
			XRequestId: None,
		}, {
			Authorization: None,
		}, {
			UserPw: "qwe",
		}, {
			UserId: "qwe",
		}, {
			UserId: "qwerqwerqwerqwerqwer",
		}, {
			UserPw: "qwerqwerqwerqwerqwer",
		}, {
			Name: "박진홍입니다",
		}, {
			Name: "응",
		}, {
			PhoneNumber: "0108837834701088378347",
		}, {
			Email: "itIsNotEmailFormat",
		}, {
			Email: "itIsSoVeryTooLongEmail@naver.com",
		},
	}

	for _, form := range forms {
		form.ExpectCode = http.StatusBadRequest
		form.ExpectMessage = MessageBadRequest
		tests = append(tests, form.createTestFromForm())
	}

	for _, test := range tests {
		test.setRequestContext(req)
		test.onExpectMethods()
		_ = h.BeforeCreateAuth(ctx, req, resp)
		assert.Equalf(t, test.ExpectCode, resp.Status, "status assert error test case: %v\n", test)
		assert.Equalf(t, test.ExpectMessage, resp.Message, "message assert error test case: %v\n", test)
	}

	mockStore.AssertExpectations(t)
}
