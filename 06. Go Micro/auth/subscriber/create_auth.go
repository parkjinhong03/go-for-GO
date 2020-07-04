package subscriber

import (
	"auth/dao"
	"auth/dao/user"
	"auth/model"
	proto "auth/proto/auth"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/micro/go-micro/v2/broker"
)

func (m *auth) CreateAuth(e broker.Event) error {
	header := e.Message().Header
	msgId := header["MessageID"]
	xReqId := header["XRequestID"]

	if xReqId == "" || msgId == "" {
		return ErrorForbidden
	}

	if _, err := uuid.Parse(xReqId); err != nil {
		return ErrorForbidden
	}

	body := proto.CreateAuthMessage{}
	if err := json.Unmarshal(e.Message().Body, &body); err != nil {
		// 에러 기록
		return ErrorBadRequest
	}

	if err := m.validate.Struct(&body); err != nil {
		// 에러 기록
		return ErrorBadRequest
	}

	var ad dao.AuthDAOService
	env, ok := e.Message().Header["Env"]
	if ok && env == "Test" {
		ad = m.adc.GetTestAuthDAO(e.(*CustomEvent).mock)
	} else {
		ad = m.adc.GetDefaultAuthDAO()
	}

	if _, err := ad.InsertMessage(&model.ProcessedMessage{
		MsgId: msgId,
	}); err != nil {
		return ErrorDuplicated
	}

	var _ *model.Auth
	if _, err := ad.InsertAuth(&model.Auth{
		UserId: body.UserId,
		UserPw: body.UserPw,
		Status: user.CreatePending,
	}); err != nil {
		ad.Rollback()
		// 에러 기록
		return nil
	}
	ad.Commit()
	// 정상 처리 기록

	_ = e.Ack()
	return nil
}