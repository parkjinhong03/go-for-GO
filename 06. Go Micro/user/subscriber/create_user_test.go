package subscriber

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log"
	"testing"
	"time"
	daoUser "user/dao/user"
	"user/model"
	proto "user/proto/user"
	"user/tool/random"
)

type createUserTest struct {
	AuthId         uint32
	Name           string
	PhoneNumber    string
	Email          string
	Introduction   string
	XRequestID     string
	MessageID      string
	AfterMessageID string
	ExpectMethods  map[method]returns
	ExpectError    error
}

func (c createUserTest) createTestFromForm() (test createUserTest) {
	test = c

	if c.Name == none 		 	{ test.Name = "" } 			 else if c.Name == "" 			{ test.Name = defaultName }
	if c.PhoneNumber == none 	{ test.PhoneNumber = "" }  	 else if c.PhoneNumber == "" 	{ test.PhoneNumber = defaultPN }
	if c.Email == none		 	{ test.Email = "" } 		 else if c.Email == "" 			{ test.Email = defaultEmail }
	if c.Email == none 			{ test.Email = "" } 		 else if c.Email == "" 		  	{ test.Email = defaultEmail }
	if c.XRequestID == none 	{ test.XRequestID = "" }	 else if c.XRequestID == ""	  	{ test.XRequestID = uuid.New().String() }
	if c.MessageID == none      { test.MessageID = "" }		 else if c.MessageID == ""	    { test.MessageID = random.GenerateString(32) }
	if c.AfterMessageID == none { test.AfterMessageID = "" } else if c.AfterMessageID == ""	{ test.AfterMessageID = random.GenerateString(32) }

	if _, ok := c.ExpectMethods["InsertUser"]; ok {
		c.setUserContext(c.ExpectMethods["InsertUser"][0].(*model.User))
	}

	if _, ok := c.ExpectMethods["InsertMessage"]; ok {
		c.setProcessedMessageContext(c.ExpectMethods["InsertMessage"][0].(*model.ProcessedMessage))
	}

	return
}

func (c createUserTest) setUserContext(user *model.User) {
	user.ID = userId
	user.AuthId = uint(c.AuthId)
	user.Name = c.Name
	user.PhoneNumber = c.PhoneNumber
	user.Email = c.Email
	user.Introduction = c.Introduction
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	userId++
}

func (c createUserTest) setProcessedMessageContext(msg *model.ProcessedMessage) {
	msg.ID = psMsgId
	msg.MsgId = c.MessageID
	msg.CreatedAt = time.Now()
	msg.UpdatedAt = time.Now()
	psMsgId++
}

func (c createUserTest) setMessageContext(msg *proto.CreateUserMessage) {
	msg.AuthId = c.AuthId
	msg.Name = c.Name
	msg.PhoneNumber = c.PhoneNumber
	msg.Email = c.Email
	msg.Introduction = c.Introduction
}

func (c createUserTest) onExpectMethods() {
	for method, returns := range c.ExpectMethods {
		c.onMethod(method, returns)
	}
}

func (c createUserTest) onMethod(method method, returns returns) {
	switch method {
	case "InsertUser":
		mockStore.On("InsertUser", &model.User{
			AuthId:		  uint(c.AuthId),
			Name:         c.Name,
			PhoneNumber:  c.PhoneNumber,
			Email:        c.Email,
			Introduction: c.Introduction,
		}).Return(returns...)
	case "InsertMessage":
		mockStore.On("InsertMessage", &model.ProcessedMessage{
			MsgId: c.MessageID,
		}).Return(returns...)
	case "Commit":
		mockStore.On("Commit").Return(returns...)
	case "Rollback":
		mockStore.On("Rollback").Return(returns...)
	case "Ack":
		mockStore.On("Ack").Return(returns...)
	//case "Publish":
	//	header := c.generateAfterMsgHeader()
	//
	//	var id uint32
	//	if _, ok := c.ExpectMethods["Insert"]; ok {
	//		id = uint32(c.ExpectMethods["Insert"][0].(*model.Auth).ID)
	//	}
	//
	//	msg := userProto.CreateUserMessage{
	//		AuthId:       id,
	//		Name:         c.Name,
	//		PhoneNumber:  c.PhoneNumber,
	//		Email:        c.Email,
	//		Introduction: c.Introduction,
	//	}
	//	body, err := json.Marshal(msg)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	mockStore.On("Publish", subscriber.CreateUserEventTopic, &broker.Message{
	//		Header: header,
	//		Body:   body,
	//	}).Return(returns...)

	// 분산 추적 관련 메서드 추가
	default:
		panic(fmt.Sprintf("%s method cannot be on booked\n", method))
	}
}

func (c createUserTest) generateMsgHeader() (header map[string]string) {
	header = make(map[string]string)
	header["XRequestID"] = c.XRequestID
	header["MessageID"] = c.MessageID
	header["AfterMessageID"] = c.AfterMessageID
	header["Env"] = "Test"

	return
}

func TestCreateUserValidMessage(t *testing.T) {
	setUpEnv()
	msg := &proto.CreateUserMessage{}
	var tests []createUserTest

	forms := []createUserTest{
		{
			AuthId: 1,
			ExpectMethods: map[method]returns{
				"InsertMessage": {&model.ProcessedMessage{}, nil},
				"InsertUser":    {&model.User{}, nil},
				"Commit":        {&gorm.DB{}},
				"Ack":           {nil},
			},
			ExpectError: nil,
		}, {
			AuthId: 2,
			ExpectMethods: map[method]returns{
				"InsertMessage": {&model.ProcessedMessage{}, nil},
				"InsertUser":    {&model.User{}, nil},
				"Commit":        {&gorm.DB{}},
				"Ack":           {errors.New("unable to ack this message")},
			},
			ExpectError: nil,
		}, {
			AuthId: 3,
			ExpectMethods: map[method]returns{
				"InsertMessage": {&model.ProcessedMessage{}, nil},
				"InsertUser":    {&model.User{}, errors.New("unable to insert user")},
				"Rollback":      {&gorm.DB{}},
			},
			ExpectError: nil,
		}, {
			AuthId: 4,
			ExpectMethods: map[method]returns{
				"InsertMessage": {&model.ProcessedMessage{}, errors.New("can't read this table")},
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

		err = h.CreateUser(event)
		assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v)\n", test)

		mockStore.AssertExpectations(t)
	}
}

func TestCreateUserUnmarshalErrorMessage(t *testing.T) {
	setUpEnv()
	msg := &proto.CreateUserMessage{}
	var tests []createUserTest

	forms := []createUserTest{{ExpectError: ErrorBadRequest}}

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

		err := h.CreateUser(event)
		assert.Equalf(t, test.ExpectError, err, "error assert error (test case: %v)\n", test)
		mockStore.AssertExpectations(t)
	}
}

func TestCreateAuthDuplicatedMessage(t *testing.T) {
	setUpEnv()
	msg := &proto.CreateUserMessage{}
	var tests []createUserTest

	forms := []createUserTest{
		{
			ExpectMethods: map[method]returns{
				"InsertMessage": {&model.ProcessedMessage{}, daoUser.MessageDuplicatedError},
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

		err = h.CreateUser(event)
		assert.Equalf(t, test.ExpectError, err, "error assert error (test case: %v)\n", test)

		mockStore.AssertExpectations(t)
	}
}

func TestCreateAuthForbiddenMessage(t *testing.T) {
	setUpEnv()
	msg := &proto.CreateUserMessage{}
	var tests []createUserTest

	forms := []createUserTest{
		{
			XRequestID: none,
		}, {
			XRequestID: "ThisIsInvalidXRequestIDString",
		}, {
			MessageID: none,
		}, {
			AfterMessageID: none,
		}, {
			MessageID: "LengthOfThisMessageIDIsNotThirtyTwo",
		}, {
			AfterMessageID: "LengthOfThisAfterMessageIDIsNotThirtyTwo",
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

		err := h.CreateUser(event)
		assert.Equalf(t, test.ExpectError, err, "error assert error (test case: %v)\n", test)

		mockStore.AssertExpectations(t)
	}
}

func TestCreateAuthBadRequestMessage(t *testing.T) {
	setUpEnv()
	msg := &proto.CreateUserMessage{}
	var tests []createUserTest

	forms := []createUserTest{
		{
			Name: none,
		}, {
			Name: "박진홍입니다",
		}, {
			Name: "응",
		}, {
			PhoneNumber: none,
		}, {
			PhoneNumber: "0108837834701088378347",
		}, {
			Email: none,
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

		err := h.CreateUser(event)
		assert.Equalf(t, test.ExpectError, err, "error assert error (test case: %v)\n", test)

		mockStore.AssertExpectations(t)
	}
}