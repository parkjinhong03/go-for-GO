package subscriber

import (
	"auth/dao"
	"auth/dao/user"
	"auth/model"
	proto "auth/proto/auth"
	"encoding/json"
	"github.com/micro/go-micro/v2/broker"
)

func (m *auth) CreateAuth(e broker.Event) error {
	header := e.Message().Header
	if header["XRequestId"] == "" || header["MessageId"] == "" || len(header["MessageId"]) != 32 {
		return ErrorBadRequest
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
		MsgId: header["MessageId"],
	}); err != nil {
		return ErrorDuplicatedMessage
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