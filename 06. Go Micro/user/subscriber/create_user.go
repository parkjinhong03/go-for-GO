package subscriber

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/micro/go-micro/v2/broker"
	"user/dao"
	userDAO "user/dao/user"
	"user/model"
	proto "user/proto/user"
	"user/tool/random"
)

func (u *user) CreateUser(event broker.Event) error {
	header := event.Message().Header
	xReqId := header["XRequestID"]
	msgId := header["MessageID"]

	if xReqId == "" || len(msgId) != 32 {
		return ErrorForbidden
	}

	if _, err := uuid.Parse(xReqId); err != nil {
		return ErrorForbidden
	}

	recvMsg := proto.CreateUserMessage{}
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

	if _, err := ud.InsertMessage(&model.ProcessedMessage{
		MsgId: msgId,
	}); err == userDAO.MessageDuplicatedError {
		return ErrorMsgDuplicated
	} else if err != nil {
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
		return nil
	}
	ud.Commit()

	// publish 추가

	if err := event.Ack(); err != nil {
		//
		return nil
	}

	return nil
}