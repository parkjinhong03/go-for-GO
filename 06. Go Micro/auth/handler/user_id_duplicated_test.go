package handler

import (
	proto "auth/proto/auth"
	"auth/tool/jwt"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"testing"
	"time"
)

type userIdDuplicatedTest struct {
	UserId              string
	Authorization       string
	XRequestId          string
	ExpectCode          uint32
	ExpectMessage       string
	ExpectMethods       map[method]returns
	ExpectAuthorization string
}

func (c userIdDuplicatedTest) createTestFromForm() (test userIdDuplicatedTest) {
	test = c

	if c.UserId == None 		{ test.UserId = "" } 		else if c.UserId == "" 		  { test.UserId = DefaultUserId }
	if c.Authorization == None 	{ test.Authorization = "" } else if c.Authorization == "" { test.Authorization = "" }
	if c.XRequestId == None		{ test.XRequestId = "" }	else if c.XRequestId == ""	  { test.XRequestId = uuid.New().String() }

	return
}

func (c userIdDuplicatedTest) setRequestContext(req *proto.UserIdDuplicatedRequest) {
	req.UserId = c.UserId
	req.Authorization = c.Authorization
	req.XRequestID = c.XRequestId
}

func (c userIdDuplicatedTest) onExpectMethods() {
	for name, returns := range c.ExpectMethods {
		c.onMethod(name, returns)
	}
}

func (c userIdDuplicatedTest) onMethod(method method, returns returns) {
	switch method {
	case "CheckIfUserIdExist":
		mockStore.On("CheckIfUserIdExist", c.UserId).Return(returns...)
	default:
		log.Fatalf("%s method cannot be on booked\n", method)
	}
	return
}

func TestUserIdDuplicatedStatusOK(t *testing.T) {
	setUpEnv()
	req := &proto.UserIdDuplicatedRequest{}
	resp := &proto.UserIdDuplicatedResponse{}
	var tests []userIdDuplicatedTest

	forms := []userIdDuplicatedTest{
		{
			UserId: "TestId1",
			ExpectMethods: map[method]returns{
				"CheckIfUserIdExist": {false, nil},
			},
			ExpectCode: http.StatusOK,
			ExpectMessage: MessageUserIdNotDuplicated,
			ExpectAuthorization: jwt.GenerateDuplicateCertJWTNoReturnErr("TestId1", "", time.Hour),
		}, {
			UserId: "TestId1",
			Authorization: jwt.GenerateDuplicateCertJWTNoReturnErr("", "jinhong0719@naver.com", time.Hour),
			ExpectMethods: map[method]returns{
				"CheckIfUserIdExist": {false, nil},
			},
			ExpectCode: http.StatusOK,
			ExpectMessage: MessageUserIdNotDuplicated,
			ExpectAuthorization: jwt.GenerateDuplicateCertJWTNoReturnErr("TestId1", "jinhong0719@naver.com", time.Hour),
		}, {
			UserId: "TestId1",
			Authorization: jwt.GenerateDuplicateCertJWTNoReturnErr("TestId2", "", time.Hour),
			ExpectMethods: map[method]returns{
				"CheckIfUserIdExist": {false, nil},
			},
			ExpectCode: http.StatusOK,
			ExpectMessage: MessageUserIdNotDuplicated,
			ExpectAuthorization: jwt.GenerateDuplicateCertJWTNoReturnErr("TestId1", "", time.Hour),
		},
	}

	for _, form := range forms {
		tests = append(tests, form.createTestFromForm())
	}

	for _, test := range tests {
		test.setRequestContext(req)
		test.onExpectMethods()
		_ = h.UserIdDuplicated(ctx, req, resp)
		assert.Equalf(t, test.ExpectCode, resp.Status, "status assertion error (test case: %v)\n", test)
		assert.Equalf(t, test.ExpectMessage, resp.Message, "message assertion error (test case: %v)\n", test)
		assert.Equalf(t, test.ExpectAuthorization, resp.Authorization, "authorization assertion error (test case: %v)\n", test)
	}
}

func TestUserIdDuplicatedDuplicateError(t *testing.T) {
	setUpEnv()
	req := &proto.UserIdDuplicatedRequest{}
	resp := &proto.UserIdDuplicatedResponse{}
	var tests []userIdDuplicatedTest

	forms := []userIdDuplicatedTest{
		{
			UserId: "TestId1",
			ExpectMethods: map[method]returns{
				"CheckIfUserIdExist": {true, nil},
			},
			ExpectCode: StatusUserIdDuplicate,
			ExpectMessage: MessageUserIdDuplicate,
		}, {
			UserId: "TestId1",
			Authorization: jwt.GenerateDuplicateCertJWTNoReturnErr("TestId1", "jinhong0719@naver.com", time.Hour),
			ExpectMethods: map[method]returns{
				"CheckIfUserIdExist": {true, nil},
			},
			ExpectCode: StatusUserIdDuplicate,
			ExpectMessage: MessageUserIdDuplicate,
		},
	}

	for _, form := range forms {
		tests = append(tests, form.createTestFromForm())
	}

	for _, test := range tests {
		test.setRequestContext(req)
		test.onExpectMethods()
		_ = h.UserIdDuplicated(ctx, req, resp)
		assert.Equalf(t, test.ExpectCode, resp.Status, "status assertion error (test case: %v)\n", test)
		assert.Equalf(t, test.ExpectMessage, resp.Message, "message assertion error (test case: %v)\n", test)
	}
}

func TestUserIdDuplicatedForbidden(t *testing.T) {
	setUpEnv()
	req := &proto.UserIdDuplicatedRequest{}
	resp := &proto.UserIdDuplicatedResponse{}
	var tests []userIdDuplicatedTest

	forms := []userIdDuplicatedTest{
		{
			UserId:        "TestId1",
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
		_ = h.UserIdDuplicated(ctx, req, resp)
		assert.Equalf(t, test.ExpectCode, resp.Status, "status assertion error (test case: %v)\n", test)
	}
}

func TestUserIdDuplicatedBadRequest(t *testing.T) {
	setUpEnv()
	req := &proto.UserIdDuplicatedRequest{}
	resp := &proto.UserIdDuplicatedResponse{}
	var tests []userIdDuplicatedTest

	forms := []userIdDuplicatedTest{
		{
			UserId: None,
		}, {
			XRequestId: None,
		}, {
			UserId: "400",
		}, {
			UserId: "thisUserIdIsTooLongMaybe400?",
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
		_ = h.UserIdDuplicated(ctx, req, resp)
		assert.Equalf(t, test.ExpectCode, resp.Status, "status assertion error (test case: %v)\n", test)
		assert.Equalf(t, test.ExpectMessage, resp.Message, "message assertion error (test case: %v)\n", test)
	}
}

func TestUserIdDuplicatedServerError(t *testing.T) {
	setUpEnv()
	req := &proto.UserIdDuplicatedRequest{}
	resp := &proto.UserIdDuplicatedResponse{}
	var tests []userIdDuplicatedTest

	forms := []userIdDuplicatedTest{
		{
			UserId: "TestId1",
			ExpectMethods: map[method]returns{
				"CheckIfUserIdExist": {true, errors.New("")},
			},
			ExpectCode: http.StatusInternalServerError,
		},
	}

	for _, form := range forms {
		form.ExpectCode = http.StatusInternalServerError
		form.ExpectMessage = MessageBadRequest
		tests = append(tests, form.createTestFromForm())
	}

	for _, test := range tests {
		test.setRequestContext(req)
		test.onExpectMethods()
		_ = h.UserIdDuplicated(ctx, req, resp)
		assert.Equalf(t, test.ExpectCode, resp.Status, "status assertion error (test case: %v)\n", test)
	}
}