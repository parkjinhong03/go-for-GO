package subscriber

import (
	"auth/dao/auth"
	"auth/model"
	authProto "auth/proto/golang/auth"
	userProto "auth/proto/golang/user"
	"auth/tool/random"
	topic "auth/topic/golang"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/micro/go-micro/v2/broker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log"
	"testing"
	"time"
)

type method string
type returns []interface{}

type createAuthTest struct {
	UserId           string
	UserPw           string
	Name             string
	PhoneNumber      string
	Email            string
	Introduction     string
	XRequestId       string
	MessageId		 string
	AfterMessageId	 string
	ExpectMethods    map[method]returns
	ExpectError  	 error
}

func (c createAuthTest) createTestFromForm() (test createAuthTest) {
	test = c

	if c.UserId == none 		{ test.UserId = "" } 		 else if c.UserId == "" 		{ test.UserId = defaultUserId }
	if c.UserPw == none 		{ test.UserPw = "" } 		 else if c.UserPw == "" 		{ test.UserPw = defaultUserPw }
	if c.Name == none 		 	{ test.Name = "" } 			 else if c.Name == "" 		  	{ test.Name = defaultName }
	if c.PhoneNumber == none  	{ test.PhoneNumber = "" }	 else if c.PhoneNumber == ""   	{ test.PhoneNumber = defaultPN }
	if c.Introduction == none	{ test.Introduction = "" } 	 else if c.Introduction == ""  	{ test.Introduction = "" }
	if c.Email == none 			{ test.Email = "" } 		 else if c.Email == "" 		  	{ test.Email = defaultEmail }
	if c.XRequestId == none 	{ test.XRequestId = "" }	 else if c.XRequestId == ""	  	{ test.XRequestId = uuid.New().String() }
	if c.MessageId == none      { test.MessageId = "" }		 else if c.MessageId == ""	    { test.MessageId = random.GenerateString(32) }
	if c.AfterMessageId == none { test.AfterMessageId = "" } else if c.AfterMessageId == ""	{ test.AfterMessageId = random.GenerateString(32) }

	if _, ok := c.ExpectMethods["InsertAuth"]; ok {
		test.setAuthContext(c.ExpectMethods["InsertAuth"][0].(*model.Auth))
	}

	if _, ok := c.ExpectMethods["InsertMessage"]; ok {
		test.setProcessedMessageContext(c.ExpectMethods["InsertMessage"][0].(*model.ProcessedMessage))
	}

	return
}

func (c createAuthTest) setAuthContext(auth *model.Auth) {
	auth.ID = authId
	auth.UserId = c.UserId
	auth.UserPw	= c.UserPw
	auth.Status = user.CreatePending
	auth.CreatedAt = time.Now()
	auth.UpdatedAt = time.Now()
	authId++
}

func (c createAuthTest) setProcessedMessageContext(msg *model.ProcessedMessage) {
	msg.ID = psMsgId
	msg.MsgId = c.MessageId
	msg.CreatedAt = time.Now()
	msg.UpdatedAt = time.Now()
	psMsgId++
}

func (c createAuthTest) setMessageContext(req *authProto.CreateAuthMessage) {
	req.UserId = c.UserId
	req.UserPw = c.UserPw
	req.Name = c.Name
	req.Email = c.Email
	req.PhoneNumber = c.PhoneNumber
	req.Introduction = c.Introduction
}

func (c createAuthTest) onExpectMethods() {
	for method, returns := range c.ExpectMethods {
		c.onMethod(method, returns)
	}
}

func (c createAuthTest) onMethod(method method, returns returns) {
	switch method {
	case "InsertAuth":
		mockStore.On("InsertAuth", &model.Auth{
			UserId: c.UserId,
			UserPw: c.UserPw,
			Status: user.CreatePending,
		}).Return(returns...)
	case "InsertMessage":
		mockStore.On("InsertMessage",&model.ProcessedMessage{
			MsgId: c.MessageId,
		}).Return(returns...)
	case "Commit":
		mockStore.On("Commit").Return(returns...)
	case "Rollback":
		mockStore.On("Rollback").Return(returns...)
	case "Ack":
		mockStore.On("Ack").Return(returns...)
	case "Publish":
		header := c.generateAfterMsgHeader()

		var id uint32
		if _, ok := c.ExpectMethods["InsertAuth"]; ok {
			id = uint32(c.ExpectMethods["InsertAuth"][0].(*model.Auth).ID)
		}

		msg := userProto.CreateUserMessage{
			AuthId:       id,
			Name:         c.Name,
			PhoneNumber:  c.PhoneNumber,
			Email:        c.Email,
			Introduction: c.Introduction,
		}
		body, err := json.Marshal(msg)
		if err != nil { log.Fatal(err) }

		mockStore.On("Publish", topic.CreateUserEventTopic, &broker.Message{
			Header: header,
			Body:   body,
		}).Return(returns...)

	// 분산 추적 관련 메서드 추가
	default:
		panic(fmt.Sprintf("%s method cannot be on booked\n", method))
	}
}

func (c createAuthTest) generateMsgHeader() (header map[string]string) {
	header = make(map[string]string)
	header["XRequestID"] = c.XRequestId
	header["MessageID"] = c.MessageId
	header["AfterMessageID"] = c.AfterMessageId
	header["Env"] = "Test"
	return
}

func (c createAuthTest) generateAfterMsgHeader() (header map[string]string) {
	header = make(map[string]string)
	header["XRequestID"] = c.XRequestId
	header["MessageID"] = c.AfterMessageId
	return
}

func TestCreateAuthValidMessage(t *testing.T) {
	setUp()
	msg := &authProto.CreateAuthMessage{}
	var tests []createAuthTest

	forms := []createAuthTest{
		{
			UserId: "TestId1",
			ExpectMethods: map[method]returns{
				"InsertMessage": {&model.ProcessedMessage{}, nil},
				"InsertAuth":    {&model.Auth{}, nil},
				"Commit":        {&gorm.DB{}},
				"Publish":       {nil},
				"Ack":           {nil},
			},
			ExpectError: nil,
		}, {
			UserId: "TestId2",
			ExpectMethods: map[method]returns{
				"InsertMessage": {&model.ProcessedMessage{}, errors.New("can't read db")},
			},
			ExpectError: nil,
		}, {
			UserId: "TestId3",
			ExpectMethods: map[method]returns{
				"InsertMessage": {&model.ProcessedMessage{}, nil},
				"InsertAuth":    {&model.Auth{}, errors.New("user id duplicated error")},
				"Rollback":      {&gorm.DB{}},
			},
			ExpectError: nil,
		}, {
			UserId: "TestId4",
			ExpectMethods: map[method]returns{
				"InsertMessage": {&model.ProcessedMessage{}, nil},
				"InsertAuth":    {&model.Auth{}, nil},
				"Commit":        {&gorm.DB{}},
				"Publish":		 {errors.New("some error occurs while publishing message")},
			},
			ExpectError: nil,
		}, {
			UserId: "TestId5",
			ExpectMethods: map[method]returns{
				"InsertMessage": {&model.ProcessedMessage{}, nil},
				"InsertAuth":    {&model.Auth{}, nil},
				"Commit":        {&gorm.DB{}},
				"Publish":		 {nil},
				"Ack":           {errors.New("some error occurs while acknowledge message")},
			},
			ExpectError: nil,
		},
	}

	for _, form := range forms {
		tests = append(tests, form.createTestFromForm())
	}

	for _, test := range tests {
		mockStore = mock.Mock{}

		test.setMessageContext(msg)
		test.onExpectMethods()

		header := test.generateMsgHeader()
		body, err := json.Marshal(msg)
		if err != nil { log.Fatal(err) }
		event.setMessage(header, body)

		err = h.CreateAuth(event)
		assert.Equalf(t, test.ExpectError, err, "error assert error (test case: %v)\n", test)

		mockStore.AssertExpectations(t)
	}
}

func TestCreateAuthUnmarshalErrorMessage(t *testing.T) {
	setUp()
	msg := &authProto.CreateAuthMessage{}
	var tests []createAuthTest

	forms := []createAuthTest{{ExpectError: ErrorBadRequest}}

	for _, form := range forms {
		tests = append(tests, form.createTestFromForm())
	}

	for _, test := range tests {
		mockStore = mock.Mock{}

		test.setMessageContext(msg)
		test.onExpectMethods()

		header := test.generateMsgHeader()
		body := []byte("unableToUnmarshalThisByteArrToStruct")
		event.setMessage(header, body)

		err := h.CreateAuth(event)
		assert.Equalf(t, test.ExpectError, err, "error assert error (test case: %v)\n", test)
		mockStore.AssertExpectations(t)
	}
}

func TestCreateAuthDuplicatedMessage(t *testing.T) {
	setUp()
	msg := &authProto.CreateAuthMessage{}
	var tests []createAuthTest

	forms := []createAuthTest{
		{
			ExpectMethods: map[method]returns{
				"InsertMessage": {&model.ProcessedMessage{}, user.MsgIdDuplicateError},
			},
			ExpectError: ErrorMsgDuplicated,
		},
	}

	for _, form := range forms {
		tests = append(tests, form.createTestFromForm())
	}

	for _, test := range tests {
		mockStore = mock.Mock{}

		test.setMessageContext(msg)
		test.onExpectMethods()

		header := test.generateMsgHeader()
		body, err := json.Marshal(msg)
		if err != nil { log.Fatal(err) }
		event.setMessage(header, body)

		err = h.CreateAuth(event)
		assert.Equalf(t, test.ExpectError, err, "error assert error (test case: %v)\n", test)

		mockStore.AssertExpectations(t)
	}
}

func TestCreateAuthForbiddenMessage(t *testing.T) {
	setUp()
	msg := &authProto.CreateAuthMessage{}
	var tests []createAuthTest

	forms := []createAuthTest{
		{
			XRequestId: none,
		}, {
			XRequestId: "ThisIsInvalidXRequestIDString",
		}, {
			MessageId: none,
		}, {
			AfterMessageId: none,
		}, {
			MessageId: "LengthOfThisMessageIDIsNotThirtyTwo",
		}, {
			AfterMessageId: "LengthOfThisAfterMessageIDIsNotThirtyTwo",
		},
	}

	for _, form := range forms {
		form.ExpectError = ErrorForbidden
		tests = append(tests, form.createTestFromForm())
	}

	for _, test := range tests {
		mockStore = mock.Mock{}

		test.setMessageContext(msg)
		test.onExpectMethods()

		header := test.generateMsgHeader()
		body, _ := json.Marshal(msg)
		event.setMessage(header, body)

		err := h.CreateAuth(event)
		assert.Equalf(t, test.ExpectError, err, "error assert error (test case: %v)\n", test)

		mockStore.AssertExpectations(t)
	}
}

func TestCreateAuthBadRequestMessage(t *testing.T) {
	setUp()
	msg := &authProto.CreateAuthMessage{}
	var tests []createAuthTest

	forms := []createAuthTest{
		{
			UserId: none,
		}, {
			UserPw: none,
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
		form.ExpectError = ErrorBadRequest
		tests = append(tests, form.createTestFromForm())
	}

	for _, test := range tests {
		mockStore = mock.Mock{}

		test.setMessageContext(msg)
		test.onExpectMethods()

		header := test.generateMsgHeader()
		body, _ := json.Marshal(msg)
		event.setMessage(header, body)

		err := h.CreateAuth(event)
		assert.Equalf(t, test.ExpectError, err, "error assert error (test case: %v)\n", test)

		mockStore.AssertExpectations(t)
	}
}