package handler

import (
	"auth/dao"
	authProto "auth/proto/golang/auth"
	"auth/tool/jwt"
	"context"
	"github.com/google/uuid"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/opentracing/opentracing-go"
	opentracinglog "github.com/opentracing/opentracing-go/log"
	"github.com/stretchr/testify/mock"
	"github.com/uber/jaeger-client-go"
	"net/http"
	"time"
)

func (e *auth) UserIdDuplicated(ctx context.Context, req *authProto.UserIdDuplicatedRequest, rsp *authProto.UserIdDuplicatedResponse) (_ error) {
	if err := e.validate.Struct(req); err != nil {
		rsp.SetStatusAndMsg(http.StatusProxyAuthRequired, err.Error())
		return
	}
	var md metadata.Metadata
	var ok bool
	if md, ok = metadata.FromContext(ctx); !ok || md == nil {
		rsp.SetStatusAndMsg(http.StatusProxyAuthRequired, MessageUnableGetMetadata)
		return
	}
	var xId string
	if xId, ok = md.Get("X-Request-Id"); !ok {
		rsp.SetStatusAndMsg(http.StatusProxyAuthRequired, MessageThereIsNoXReqId)
		return
	}
	if _, err := uuid.Parse(xId); err != nil {
		rsp.SetStatusAndMsg(http.StatusProxyAuthRequired, err.Error())
		return
	}

	// api gateway에서 옳바르지 않은 JWT 필터링 기능 추가 필요 (proto도 변경...............해야된다니)
	var email string
	if ss, ok := md.Get("Unique-Authorization"); ok {
		claim, err := jwt.ParseDuplicateCertClaimFromJWT(ss)
		if err != nil { rsp.SetStatusAndMsg(http.StatusForbidden, err.Error()); return }
		email = claim.Email
	}

	scs, ok := md.Get("Span-Context")
	if !ok {
		rsp.SetStatusAndMsg(http.StatusProxyAuthRequired, MessageThereIsNoXReqId)
		return
	}
	sc, err := jaeger.ContextFromString(scs)
	if err != nil {
		rsp.SetStatusAndMsg(http.StatusProxyAuthRequired, err.Error())
		return
	}

	var ad dao.AuthDAOService
	switch ctx.Value("env") {
	case "test":
		mockStore := ctx.Value("mockStore").(*mock.Mock)
		ad = e.adc.GetTestAuthDAO(mockStore)
	default:
		ad = e.adc.GetDefaultAuthDAO()
	}

	dsp := e.tracer.StartSpan("CheckIfUserIdExist", opentracing.ChildOf(sc))
	dsp.SetTag("X-Request-Id", xId)
	exist, err := ad.CheckIfUserIdExist(req.UserId)
	dsp.LogFields(opentracinglog.Bool("exist", exist), opentracinglog.Error(err))
	dsp.Finish()

	if err != nil {
		rsp.SetStatusAndMsg(http.StatusInternalServerError, err.Error())
		return
	}

	if exist {
		rsp.SetStatusAndMsg(StatusUserIdDuplicate, MessageUserIdDuplicate)
		return
	}

	ss, err := jwt.GenerateDuplicateCertJWT(req.UserId, email, time.Hour)
	if err != nil {
		rsp.SetStatusAndMsg(http.StatusInternalServerError, err.Error())
		return
	}

	rsp.SetStatusAndMsg(http.StatusOK, MessageUserIdNotDuplicated)
	rsp.Authorization = ss
	return
}