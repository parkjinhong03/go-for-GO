package handler

import (
	authProto "auth/proto/golang/auth"
	"auth/tool/jwt"
	"auth/tool/random"
	topic "auth/topic/golang"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-micro/v2/metadata"
	"net/http"
)

func (e *auth) BeforeCreateAuth(ctx context.Context, req *authProto.BeforeCreateAuthRequest, rsp *authProto.BeforeCreateAuthResponse) (_ error) {
	var err error
	if err := e.validate.Struct(req); err != nil {
		rsp.SetStatusAndMsg(http.StatusBadRequest, MessageBadRequest)
		return
	}

	var md metadata.Metadata
	var ok bool
	if md, ok = metadata.FromContext(ctx); !ok || md == nil {
		rsp.SetStatus(http.StatusForbidden)
		return
	}

	var xId string
	if xId, ok = md.Get("X-Request-Id"); !ok || xId == "" {
		rsp.SetStatus(http.StatusForbidden)
		return
	}

	if _, err := uuid.Parse(xId); err != nil {
		rsp.SetStatus(http.StatusForbidden)
		return
	}

	var ss string
	if ss, ok = md.Get("Unique-Authorization"); !ok || ss == "" {
		rsp.SetStatus(http.StatusForbidden)
		return
	}

	var claim *jwt.DuplicateCertClaim
	if claim, err = jwt.ParseDuplicateCertClaimFromJWT(ss); err != nil {
		rsp.SetStatus(http.StatusForbidden)
		return
	}

	// test 환경일 경우 context에서 MessageId 추출, 아닐 시 새로 생성
	var mId string
	switch ctx.Value("env") {
	case "test":
		mId, ok = md.Get("Message-Id")
		if !ok || len(mId) != 32 { rsp.SetStatus(http.StatusForbidden); return }
	default:
		mId = random.GenerateString(32)
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
	header["XRequestID"] = xId
	header["MessageID"] = mId

	msg := authProto.CreateAuthMessage{
		UserId:       req.UserId,
		UserPw:       req.UserPw,
		Name:         req.Name,
		PhoneNumber:  req.PhoneNumber,
		Email:        req.Email,
		Introduction: req.Introduction,
	}

	body, err := json.Marshal(msg)
	if err != nil {
		rsp.SetStatus(http.StatusInternalServerError)
	}

	if err = e.mq.Publish(topic.CreateAuthEventTopic, &broker.Message{
		Header: header,
		Body:   body,
	}); err != nil {
		rsp.SetStatus(http.StatusInternalServerError)
		return
	}

	rsp.SetStatusAndMsg(http.StatusCreated, MessageAuthCreated)
	return
}