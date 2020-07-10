package subscriber

import (
	"auth/dao"
	"auth/dao/user"
	"auth/model"
	authProto "auth/proto/golang/auth"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/micro/go-micro/v2/broker"
)

func (m *auth) ChangeAuthStatus(event broker.Event) (_ error) {
	header := event.Message().Header
	xReqId := header["XRequestID"]
	msgId := header["MessageID"]
	recvMsg := authProto.ChangeAuthStatusMessage{}

	if err := json.Unmarshal(event.Message().Body, &recvMsg); err != nil {
		return ErrorBadRequest
	}

	if err := m.validate.Struct(&recvMsg); err != nil {
		return ErrorBadRequest
	}

	if xReqId == "" || len(msgId) != 32 {
		return ErrorForbidden
	}

	if _, err := uuid.Parse(xReqId); err != nil {
		return ErrorForbidden
	}

	var ad dao.AuthDAOService
	switch header["Env"] {
	case "Test":
		mock := event.(*CustomEvent).mock
		ad = m.adc.GetTestAuthDAO(mock)
	default:
		ad = m.adc.GetDefaultAuthDAO()
	}

	if _, err := ad.InsertMessage(&model.ProcessedMessage{
		MsgId: msgId,
	}); err == user.MsgIdDuplicateError {
		return ErrorMsgDuplicated
	} else if err != nil {
		return
	}

	var status = user.Created
	if !recvMsg.Success {
		status = user.Rejected
	}

	if err := ad.UpdateStatus(uint(recvMsg.AuthId), status); err != nil {
		ad.Rollback()
		return
	}

	ad.Commit()
	if err := event.Ack(); err != nil {
		//
	}

	return
}