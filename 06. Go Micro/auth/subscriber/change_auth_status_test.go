package subscriber

import (
	"auth/dao/user"
	"auth/model"
	authProto "auth/proto/golang/auth"
	"auth/tool/random"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log"
	"testing"
	"time"
)

type changeAuthStatusTest struct {
	AuthId       uint32
	Success      bool // examples.blog.service.user.CreateUser 로 부터 받은 msg의 success 필드
	ExpectError  error
	ExpectMethod map[method]returns
	XRequestID   string
	MessageID	 string
}

func (c changeAuthStatusTest) createTestFromForm() (test changeAuthStatusTest) {
	test = c

	if c.AuthId == noneInt 	{ test.AuthId = 0 } 	 else if c.AuthId == 0 		{ test.AuthId = 1 }
	if c.XRequestID == none { test.XRequestID = "" } else if c.XRequestID == "" { test.XRequestID = uuid.New().String() }
	if c.MessageID == none 	{ test.MessageID = "" }  else if c.MessageID == "" 	{ test.MessageID = random.GenerateString(32) }

	if returns, ok := c.ExpectMethod["InsertMessage"]; ok {
		c.setProcessedMessageContext(returns[0].(*model.ProcessedMessage))
	}
	return
}

func (c changeAuthStatusTest) setProcessedMessageContext(psMsg *model.ProcessedMessage) {
	psMsg.ID = psMsgId
	psMsg.MsgId = c.MessageID
	psMsg.CreatedAt = time.Now()
	psMsg.UpdatedAt = time.Now()
	psMsgId++
}

func (c changeAuthStatusTest) setMessageContext(msg *authProto.ChangeAuthStatusMessage) {
	msg.AuthId = c.AuthId
	msg.Success = c.Success
}

func (c changeAuthStatusTest) generateMsgHeader() (header map[string]string) {
	header = make(map[string]string)
	header["XRequestID"] = c.XRequestID
	header["MessageID"] = c.MessageID
	return
}

func (c changeAuthStatusTest) onExpectMethods() {
	for method, returns := range c.ExpectMethod {
		c.onMethod(method, returns)
	}
}

func (c changeAuthStatusTest) onMethod(method method, returns returns) {
	switch method {
	case "InsertMessage":
		mockStore.On("InsertMessage", &model.ProcessedMessage{
			MsgId: c.MessageID,
		}).Return(returns...)
	case "UpdateStatus":
		var status = user.Created
		if !c.Success {
			status = user.Rejected
		}
		mockStore.On("UpdateStatus", c.AuthId, status).Return(returns...)
	case "Commit":
		mockStore.On("Commit").Return(returns...)
	case "Rollback":
		mockStore.On("Rollback").Return(returns...)
	case "Ack":
		mockStore.On("Ack").Return(returns...)
	}
}

func TestAuthChangeAuthStatusValidMessage(t *testing.T) {
	setUp()
	msg := &authProto.ChangeAuthStatusMessage{}
	var tests []changeAuthStatusTest

	forms := []changeAuthStatusTest{
		{
			AuthId:  1,
			Success: true,
			ExpectMethod: map[method]returns{
				"InsertMessage": {&model.ProcessedMessage{}, nil},
				"UpdateStatus":  {nil},
				"Commit":        {&gorm.DB{}},
				"Ack":           {nil},
			},
			ExpectError: nil,
		}, {
			AuthId:  2,
			Success: true,
			ExpectMethod: map[method]returns{
				"InsertMessage": {&model.ProcessedMessage{}, nil},
				"UpdateStatus":  {nil},
				"Commit":        {&gorm.DB{}},
				"Ack":           {errors.New("unable to ack message")},
			},
			ExpectError: nil,
		}, {
			AuthId:  3,
			Success: true,
			ExpectMethod: map[method]returns{
				"InsertMessage": {&model.ProcessedMessage{}, nil},
				"UpdateStatus":  {errors.New("unable to update auth's status")},
				"Rollback":      {&gorm.DB{}},
			},
			ExpectError: nil,
		}, {
			AuthId:  4,
			Success: true,
			ExpectMethod: map[method]returns{
				"InsertMessage": {&model.ProcessedMessage{}, errors.New("unable to read db")},
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
		err = h.ChangeAuthStatus(event)

		assert.Equalf(t, test.ExpectError, err, "error assertion error (test caseL %v)\n", test)
		mockStore.AssertExpectations(t)
	}
}

func TestAuthChangeAuthStatusUnmarshalErrorMessage(t *testing.T) {
	setUp()
	msg := &authProto.ChangeAuthStatusMessage{}
	var tests []changeAuthStatusTest

	forms := []changeAuthStatusTest{{ ExpectError: ErrorBadRequest }}

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
		err := h.ChangeAuthStatus(event)

		assert.Equalf(t, test.ExpectError, err, "error assertion error (test caseL %v)\n", test)
		mockStore.AssertExpectations(t)
	}
}