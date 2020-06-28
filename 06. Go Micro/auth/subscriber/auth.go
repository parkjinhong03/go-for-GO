package subscriber

import (
	"auth/dao"
	"auth/dao/user"
	"auth/model"
	proto "auth/proto/auth"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/micro/go-micro/v2/broker"
)

type msgHandler struct {
	adc      *dao.AuthDAOCreator
	validate *validator.Validate
}

func NewMsgHandler(adc *dao.AuthDAOCreator, validate *validator.Validate) *msgHandler {
	return &msgHandler{
		adc:      adc,
		validate: validate,
	}
}

func (m *msgHandler) CreateAuth(e broker.Event) (err error) {
	body := proto.CreateAuthMessage{}
	if err = json.Unmarshal(e.Message().Body, &body); err != nil {
		// 에러 기록
		return
	}

	if err = m.validate.Struct(&body); err != nil {
		// 에러 기록
		return
	}

	// 중복 메시지 처리 로직 추가

	var ad dao.AuthDAOService
	env, ok := e.Message().Header["Env"]
	if ok && env == "Test" {
		ad = m.adc.GetTestAuthDAO(e.(*CustomEvent).mock)
	} else {
		ad = m.adc.GetDefaultAuthDAO()
	}

	var _ *model.Auth
	if _, err = ad.Insert(&model.Auth{
		UserId: body.UserId,
		UserPw: body.UserPw,
		Status: user.CreatePending,
	}); err != nil {
		ad.Rollback()
		// 에러 기록
		return nil
	}

	if err := e.Ack(); err != nil {
		ad.Rollback()
		// 에러 기록
	} else {
		// 정상 처리 기록
	}

	ad.Commit()
	return nil
}
