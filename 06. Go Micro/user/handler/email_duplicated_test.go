package handler

import (
	"errors"
	"github.com/google/uuid"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"testing"
	"time"
	proto "user/proto/user"
	"user/tool/jwt"
)

type emailDuplicatedTest struct {
	Email               string
	Authorization       string
	XRequestId          string
	ExpectCode          uint32
	ExpectMessage       string
	ExpectMethods       map[method]returns
	ExpectAuthorization string
}

func (e emailDuplicatedTest) createTestFromForm() (test emailDuplicatedTest) {
	test = e

	if e.Email == none 			{ test.Email = "" } 	 	else if e.Email == "" 		  { test.Email = defaultEmail }
	if e.Authorization == none  { test.Authorization = "" } else if e.Authorization == "" { test.Authorization = "" }
	if e.XRequestId == none 	{ test.XRequestId = "" } 	else if e.XRequestId == "" 	  { test.XRequestId = uuid.New().String() }

	return
}

func (e emailDuplicatedTest) setRequestContext(req *proto.EmailDuplicatedRequest) {
	req.Email = e.Email
	ctx = metadata.Set(ctx, "XRequestID", e.XRequestId)
	ctx = metadata.Set(ctx, "Authorization", e.Authorization)
	return
}

func (e emailDuplicatedTest) onExpectMethods() {
	for method, returns := range e.ExpectMethods {
		e.onMethod(method, returns)
	}
	return
}

func (e emailDuplicatedTest) onMethod(method method, returns returns) {
	switch method {
	case "CheckIfEmailExist":
		mockStore.On("CheckIfEmailExist", e.Email).Return(returns...)
	default:
		log.Fatalf("%s method cannot be on booked\n", method)
	}
	return
}

func TestEmailDuplicatedStatusOK(t *testing.T) {
	setUpEnv()
	req := &proto.EmailDuplicatedRequest{}
	resp := &proto.EmailDuplicatedResponse{}
	var tests []emailDuplicatedTest

	forms := []emailDuplicatedTest{
		{
			Email: "jinhong0719@naver.com",
			ExpectMethods: map[method]returns{
				"CheckIfEmailExist": {false, nil},
			},
			ExpectCode: http.StatusOK,
			ExpectMessage: MessageEmailNotDuplicated,
			ExpectAuthorization: jwt.GenerateDuplicateCertJWTNoReturnErr("", "jinhong0719@naver.com", time.Hour),
		}, {
			Email: "jinhong0719@naver.com",
			Authorization: jwt.GenerateDuplicateCertJWTNoReturnErr("TestId", "", time.Hour),
			ExpectMethods: map[method]returns{
				"CheckIfEmailExist": {false, nil},
			},
			ExpectCode: http.StatusOK,
			ExpectMessage: MessageEmailNotDuplicated,
			ExpectAuthorization: jwt.GenerateDuplicateCertJWTNoReturnErr("TestId", "jinhong0719@naver.com", time.Hour),
		}, {
			Email: "jinhong0719@naver.com",
			Authorization: jwt.GenerateDuplicateCertJWTNoReturnErr("", "jinhong0719@hanmail.net", time.Hour),
			ExpectMethods: map[method]returns{
				"CheckIfEmailExist": {false, nil},
			},
			ExpectCode: http.StatusOK,
			ExpectMessage: MessageEmailNotDuplicated,
			ExpectAuthorization: jwt.GenerateDuplicateCertJWTNoReturnErr("", "jinhong0719@naver.com", time.Hour),
		},
	}

	for _, form := range forms {
		tests = append(tests, form.createTestFromForm())
	}

	for _, test := range tests {
		test.setRequestContext(req)
		test.onExpectMethods()
		_ = h.EmailDuplicated(ctx, req, resp)
		assert.Equalf(t, int(test.ExpectCode), int(resp.Status), "status assertion error (test case: %v)\n", test)
		assert.Equalf(t, test.ExpectMessage, resp.Message, "message assertion error (test case: %v)\n", test)
		assert.Equalf(t, test.ExpectAuthorization, resp.Authorization, "authorization assertion error (test case: %v)\n", test)
	}
}

func TestEmailDuplicatedDuplicateError(t *testing.T) {
	setUpEnv()
	req := &proto.EmailDuplicatedRequest{}
	resp := &proto.EmailDuplicatedResponse{}
	var tests []emailDuplicatedTest

	forms := []emailDuplicatedTest{
		{
			ExpectMethods: map[method]returns{
				"CheckIfEmailExist": {true, nil},
			},
			ExpectCode: StatusEmailDuplicated,
			ExpectMessage: MessageEmailDuplicated,
		}, {
			Authorization: jwt.GenerateDuplicateCertJWTNoReturnErr("TestId1", "jinhong0719@naver.com", time.Hour),
			ExpectMethods: map[method]returns{
				"CheckIfEmailExist": {true, nil},
			},
			ExpectCode: StatusEmailDuplicated,
			ExpectMessage: MessageEmailDuplicated,
		},
	}

	for _, form := range forms {
		tests = append(tests, form.createTestFromForm())
	}

	for _, test := range tests {
		test.setRequestContext(req)
		test.onExpectMethods()
		_ = h.EmailDuplicated(ctx, req, resp)
		assert.Equalf(t, int(test.ExpectCode), int(resp.Status), "status assertion error (test case: %v)\n", test)
		assert.Equalf(t, test.ExpectMessage, resp.Message, "message assertion error (test case: %v)\n", test)
	}
}

func TestEmailDuplicatedForbidden(t *testing.T) {
	setUpEnv()
	req := &proto.EmailDuplicatedRequest{}
	resp := &proto.EmailDuplicatedResponse{}
	var tests []emailDuplicatedTest

	forms := []emailDuplicatedTest{
		{
			Email:         "jinhong07191@naver.com",
			Authorization: "ThisIsInvalidAuthorizationString",
		}, {
			Email:      "jinhong07192@naver.com",
			XRequestId: none,
		}, {
			Email:      "jinhong07193@naver.com",
			XRequestId: "ThisIsInvalidXRequestIDString",
		},
	}

	for _, form := range forms {
		form.ExpectCode = http.StatusForbidden
		tests = append(tests, form.createTestFromForm())
	}

	for _, test := range tests {
		test.setRequestContext(req)
		test.onExpectMethods()
		_ = h.EmailDuplicated(ctx, req, resp)
		assert.Equalf(t, int(test.ExpectCode), int(resp.Status), "status assertion error (test case: %v)\n", test)
	}
}

func TestEmailDuplicatedBadRequest(t *testing.T) {
	setUpEnv()
	req := &proto.EmailDuplicatedRequest{}
	resp := &proto.EmailDuplicatedResponse{}
	var tests []emailDuplicatedTest

	forms := []emailDuplicatedTest{
		{
			Email: none,
		}, {
			Email: "thisEmailIsTooLongMaybe400?@naver.com",
		}, {
			Email: "thisEmailIsInvalid",
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
		_ = h.EmailDuplicated(ctx, req, resp)
		assert.Equalf(t, int(test.ExpectCode), int(resp.Status), "status assertion error (test case: %v)\n", test)
		assert.Equalf(t, test.ExpectMessage, resp.Message, "message assertion error (test case: %v)\n", test)
	}
}

func TestEmailDuplicatedServerError(t *testing.T) {
	setUpEnv()
	req := &proto.EmailDuplicatedRequest{}
	resp := &proto.EmailDuplicatedResponse{}
	var tests []emailDuplicatedTest

	forms := []emailDuplicatedTest{
		{
			ExpectMethods: map[method]returns{
				"CheckIfEmailExist": {true, errors.New("some db error")},
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
		_ = h.EmailDuplicated(ctx, req, resp)
		assert.Equalf(t, int(test.ExpectCode), int(resp.Status), "status assertion error (test case: %v)\n", test)
	}
}