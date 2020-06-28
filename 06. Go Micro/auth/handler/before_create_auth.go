package handler

import (
	proto "auth/proto/auth"
	"auth/subscriber"
	"auth/tool/jwt"
	"auth/tool/random"
	"context"
	"encoding/json"
	"github.com/micro/go-micro/v2/broker"
	"net/http"
)

func (e *auth) BeforeCreateAuth(ctx context.Context, req *proto.BeforeCreateAuthRequest, rsp *proto.BeforeCreateAuthResponse) (_ error) {
	if err := e.validate.Struct(req); err != nil {
		rsp.SetStatusAndMsg(http.StatusBadRequest, MessageBadRequest)
		return
	}

	// test 환경 시 context에서 MessageId 추출, 아닐 시 새로 생성
	var mId string
	switch ctx.Value("env") {
	case "test":
		mId = ctx.Value("MessageId").(string)
	default:
		mId = random.GenerateString(32)
	}

	claim, err := jwt.ParseDuplicateCertClaimFromJWT(req.Authorization)
	if err != nil {
		rsp.SetStatus(http.StatusForbidden)
		return
	}

	if claim.UserId != req.UserId {
		rsp.SetStatusAndMsg(StatusUserIdDuplicate, MessageUserIdDuplicate)
		return
	}

	if claim.Email != req.Email {
		rsp.SetStatusAndMsg(StatusEmailDuplicate, MessageEmailDuplicate)
		return
	}

	header := make(map[string]string)
	header["XRequestId"] = req.XRequestID
	header["MessageId"] = mId

	body := proto.CreateAuthMessage{
		UserId:       req.UserId,
		UserPw:       req.UserPw,
		Name:         req.Name,
		PhoneNumber:  req.PhoneNumber,
		Email:        req.Email,
		Introduction: req.Introduction,
	}

	b, err := json.Marshal(body)
	if err != nil {
		rsp.SetStatus(http.StatusInternalServerError)
	}

	if err = e.mq.Publish(subscriber.CreateAuthEventTopic, &broker.Message{
		Header: header,
		Body:   b,
	}); err != nil {
		rsp.SetStatus(http.StatusInternalServerError)
		return
	}

	rsp.SetStatusAndMsg(http.StatusCreated, MessageAuthCreated)
	return
}