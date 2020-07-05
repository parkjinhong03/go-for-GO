package subscriber

import (
	"fmt"
	"github.com/google/uuid"
	"time"
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