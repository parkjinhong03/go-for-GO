package subscriber

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/micro/go-micro/v2/broker"
	"user/dao"
	userDAO "user/dao/user"
	"user/model"
	authProto "user/proto/golang/auth"
	userProto "user/proto/golang/user"
	"user/tool/random"
	topic "user/topic/golang"
)

func (u *User) CreateUser(event broker.Event) error {
	header := event.Message().Header
	xReqId := header["XRequestID"]
	msgId := header["MessageID"]

	if xReqId == "" || len(msgId) != 32 {
		return ErrorForbidden
	}

	if _, err := uuid.Parse(xReqId); err != nil {
		return ErrorForbidden
	}

	recvMsg := userProto.CreateUserMessage{}
	if err := json.Unmarshal(event.Message().Body, &recvMsg); err != nil {
		return ErrorBadRequest
	}

	if err := u.validate.Struct(&recvMsg); err != nil {
		return ErrorBadRequest
	}

	var ud dao.UserDAOService
	var aftMsgId string
	switch header["Env"] {
	case "Test":
		mock := event.(*CustomEvent).mock
		ud = u.udc.GetTestUserDAO(mock)
		aftMsgId = header["AfterMessageID"]
		if len(aftMsgId) != 32 { return ErrorForbidden }
	default:
		ud = u.udc.GetDefaultUserDAO()
		aftMsgId = random.GenerateString(32)
	}

	header = make(map[string]string)
	header["XRequestID"] = xReqId
	header["MessageID"] = aftMsgId
	brkMsg := &broker.Message{Header: header}
	sendMsg := authProto.ChangeAuthStatusMessage{AuthId: recvMsg.AuthId}

	if _, err := ud.InsertMessage(&model.ProcessedMessage{
		MsgId: msgId,
	}); err == userDAO.MessageDuplicatedError {
		return ErrorMsgDuplicated
	} else if err != nil {
		sendMsg.Success = false
		brkMsg.Body, _ = json.Marshal(sendMsg)
		_ = u.mq.Publish(topic.ChangeAuthStatusEventTopic, brkMsg)
		return nil
	}

	_, err := ud.InsertUser(&model.User{
		AuthId:       uint(recvMsg.AuthId),
		Name:         recvMsg.Name,
		PhoneNumber:  recvMsg.PhoneNumber,
		Email:        recvMsg.Email,
		Introduction: recvMsg.Introduction,
	})

	if err != nil {
		ud.Rollback()
		sendMsg.Success = false
		brkMsg.Body, _ = json.Marshal(sendMsg)
		_ = u.mq.Publish(topic.ChangeAuthStatusEventTopic, brkMsg)
		return nil
	}
	ud.Commit()

	sendMsg.Success = true
	brkMsg.Body, _ = json.Marshal(sendMsg)
	if err := u.mq.Publish(topic.ChangeAuthStatusEventTopic, brkMsg); err != nil {
		// 로직상 성공, 하지만 publish 에러
		return nil
	}

	if err := event.Ack(); err != nil {
		// publish 까지 성공, 하지만 ack 에러
		return nil
	}

	return nil
}