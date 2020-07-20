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
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	"net/http"
)

func (e *auth) BeforeCreateAuth(ctx context.Context, req *authProto.BeforeCreateAuthRequest, rsp *authProto.BeforeCreateAuthResponse) (_ error) {
	var err error
	if err := e.validate.Struct(req); err != nil {
		rsp.SetStatus(http.StatusProxyAuthRequired)
		return
	}
	var md metadata.Metadata
	var ok bool
	if md, ok = metadata.FromContext(ctx); !ok || md == nil {
		rsp.SetStatus(http.StatusProxyAuthRequired)
		return
	}
	var xId string
	if xId, ok = md.Get("X-Request-Id"); !ok || xId == "" {
		rsp.SetStatus(http.StatusProxyAuthRequired)
		return
	}
	if _, err := uuid.Parse(xId); err != nil {
		rsp.SetStatus(http.StatusProxyAuthRequired)
		return
	}

	// 인증 api gateway에서 처리 예정 (proto 수정 필요...)
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
	var mId = random.GenerateString(32)
	if mid, ok := md.Get("Message-Id"); ok || len(mid) == 32 {
		mId = mid
	}
	if claim.UserId != req.UserId {
		rsp.SetStatusAndMsg(StatusUserIdDuplicate, MessageUserIdDuplicate)
		return
	}
	if claim.Email != req.Email {
		rsp.SetStatusAndMsg(StatusEmailDuplicate, MessageEmailDuplicate)
		return
	}

	sps, ok := md.Get("Span-Context")
	if !ok {
		rsp.SetStatus(http.StatusProxyAuthRequired)
		return
	}
	cs, err := jaeger.ContextFromString(sps)
	if err != nil {
		rsp.SetStatus(http.StatusProxyAuthRequired)
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
	body, _ := json.Marshal(msg)

	msp := e.tracer.StartSpan("MessagePublish", opentracing.ChildOf(cs))
	msp.SetTag("X-Request-Id", xId).SetTag("Published-Message-Id", mId)
	err = e.mq.Publish(topic.CreateAuthEventTopic, &broker.Message{
		Header: header,
		Body:   body,
	})
	msp.LogFields(log.Object("message", msg), log.Object("header", header), log.Error(err))
	msp.Finish()

	if err != nil {
		rsp.SetStatus(http.StatusInternalServerError)
		return
	}
	rsp.SetStatusAndMsg(http.StatusCreated, MessageAuthCreated)
	return
}