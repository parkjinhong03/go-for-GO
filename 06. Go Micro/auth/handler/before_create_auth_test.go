package handler

import (
	authProto "auth/proto/golang/auth"
	"auth/tool/jwt"
	"auth/tool/random"
	topic "auth/topic/golang"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	MessageId     string
	ExpectCode    uint32
	ExpectMessage string
	ExpectMethods map[method]returns
}

func (c CreateAuthTest) createTestFromForm() (test CreateAuthTest) {
	test = c

	if c.UserId == None 		{ test.UserId = "" } 		else if c.UserId == "" 		  { test.UserId = DefaultUserId }
	if c.UserPw == None 		{ test.UserPw = "" } 		else if c.UserPw == "" 		  { test.UserPw = DefaultUserPw }
	if c.Name == None 		 	{ test.Name = "" } 			else if c.Name == "" 		  { test.Name = DefaultName }
	if c.PhoneNumber == None  	{ test.PhoneNumber = "" }	else if c.PhoneNumber == ""   { test.PhoneNumber = DefaultPN }
	if c.Introduction == None	{ test.Introduction = "" } 	else if c.Introduction == ""  { test.Introduction = "" }
	if c.Email == None 			{ test.Email = "" } 		else if c.Email == "" 		  { test.Email = DefaultEmail }
	if c.XRequestId == None 	{ test.XRequestId = "" }	else if c.XRequestId == ""	  { test.XRequestId = uuid.New().String() }
	if c.MessageId == None      { test.MessageId = "" }     else if c.MessageId == ""     { test.MessageId = random.GenerateString(32) }
	if c.Authorization == None  { test.Authorization = "" } else if c.Authorization == "" {
		test.Authorization = jwt.GenerateDuplicateCertJWTNoReturnErr(test.UserId, test.Email, time.Hour)
	}
	return
}

func (c CreateAuthTest) setRequestContext(req *authProto.BeforeCreateAuthRequest) {
	req.UserId = c.UserId
	req.UserPw = c.UserPw
	req.Name = c.Name
	req.Email = c.Email
	req.PhoneNumber = c.PhoneNumber
	req.Introduction = c.Introduction
	ctx = metadata.Set(ctx, "X-Request-Id", c.XRequestId)
	ctx = metadata.Set(ctx, "Unique-Authorization", c.Authorization)
	ctx = metadata.Set(ctx, "Message-Id", c.MessageId)
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
		header["XRequestID"] = c.XRequestId
		header["MessageID"]  = c.MessageId

		msg := authProto.CreateAuthMessage{
			UserId:       c.UserId,
			UserPw:       c.UserPw,
			Name:         c.Name,
			PhoneNumber:  c.PhoneNumber,
			Email:        c.Email,
			Introduction: c.Introduction,
		}
		body, err := json.Marshal(msg)
		if err != nil { log.Fatal(err) }

		mockStore.On("Publish", topic.CreateAuthEventTopic, &broker.Message{
			Header: header,
			Body:   body,
		}).Return(returns...)
	default:
		panic(fmt.Sprintf("%s method cannot be on booked\n", method))
	}
	return
}

func TestAuthCreateManySuccess(t *testing.T) {
	setUpEnv()
	req := &authProto.BeforeCreateAuthRequest{}
	resp := &authProto.BeforeCreateAuthResponse{}
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
		mockStore = mock.Mock{}
		test.setRequestContext(req)
		test.onExpectMethods()
		_ = h.BeforeCreateAuth(ctx, req, resp)
		assert.Equalf(t, int(test.ExpectCode), int(resp.Status), "status assert error test case: %v\n", test)
		assert.Equalf(t, test.ExpectMessage, resp.Message, "message assert error test case: %v\n", test)

		mockStore.AssertExpectations(t)
	}
}

func TestBeforeCreateAuthUserIdDuplicateError(t *testing.T) {
	setUpEnv()
	req := &authProto.BeforeCreateAuthRequest{}
	resp := &authProto.BeforeCreateAuthResponse{}
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
		mockStore = mock.Mock{}

		test.setRequestContext(req)
		test.onExpectMethods()
		_ = h.BeforeCreateAuth(ctx, req, resp)
		assert.Equalf(t, int(test.ExpectCode), int(resp.Status), "status assert error test case: %v\n", test)
		assert.Equalf(t, test.ExpectMessage, resp.Message, "message assert error test case: %v\n", test)

		mockStore.AssertExpectations(t)
	}
}

func TestBeforeCreateAuthForbidden(t *testing.T) {
	setUpEnv()
	req := &authProto.BeforeCreateAuthRequest{}
	resp := &authProto.BeforeCreateAuthResponse{}
	var tests []CreateAuthTest

	var forms = []CreateAuthTest{
		{
			UserId:     "testId1",
			XRequestId: None,
		}, {
			UserId:     "testId2",
			XRequestId: "ThisIsInvalidXRequestIDString",
		}, {
			UserId:        "testId3",
			Authorization: None,
		}, {
			UserId:        "testId4",
			Authorization: "ThisIsInvalidAuthorizationString",
		}, {
			UserId:    "testId5",
			MessageId: None,
		}, {
			UserId:    "testId6",
			MessageId: "LengthOfThisMessageIDIsNotThirtyTwo",
		}, {
			UserId:    "testId7",
			MessageId: None,
		},
	}

	for _, form := range forms {
		form.ExpectCode = http.StatusForbidden
		tests = append(tests, form.createTestFromForm())
	}

	for _, test := range tests {
		mockStore = mock.Mock{}

		test.setRequestContext(req)
		test.onExpectMethods()
		_ = h.BeforeCreateAuth(ctx, req, resp)
		assert.Equalf(t, int(test.ExpectCode), int(resp.Status), "status assert error test case: %v\n", test)

		mockStore.AssertExpectations(t)
	}
}

func TestBeforeAuthCreateInsertBadRequest(t *testing.T) {
	setUpEnv()
	req := &authProto.BeforeCreateAuthRequest{}
	resp := &authProto.BeforeCreateAuthResponse{}
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
		mockStore = mock.Mock{}

		test.setRequestContext(req)
		test.onExpectMethods()
		_ = h.BeforeCreateAuth(ctx, req, resp)
		assert.Equalf(t, int(test.ExpectCode), int(resp.Status), "status assert error test case: %v\n", test)
		assert.Equalf(t, test.ExpectMessage, resp.Message, "message assert error test case: %v\n", test)

		mockStore.AssertExpectations(t)
	}
}

func TestBeforeCreateAuthServerError(t *testing.T) {
	setUpEnv()
	req := &authProto.BeforeCreateAuthRequest{}
	resp := &authProto.BeforeCreateAuthResponse{}
	var tests []CreateAuthTest

	var forms = []CreateAuthTest{
		{
			UserId:        "testId1",
			Email:         "richimous0719@naver.com",
			Authorization: jwt.GenerateDuplicateCertJWTNoReturnErr("testId1", "richimous0719@naver.com", time.Hour),
			ExpectMethods: map[method]returns{
				"Publish": {errors.New("")},
			},
			ExpectCode:    http.StatusInternalServerError,
		},
	}

	for _, form := range forms {
		tests = append(tests, form.createTestFromForm())
	}

	for _, test := range tests {
		mockStore = mock.Mock{}

		test.setRequestContext(req)
		test.onExpectMethods()
		_ = h.BeforeCreateAuth(ctx, req, resp)
		assert.Equalf(t, int(test.ExpectCode), int(resp.Status), "status assert error test case: %v\n", test)

		mockStore.AssertExpectations(t)
	}
}