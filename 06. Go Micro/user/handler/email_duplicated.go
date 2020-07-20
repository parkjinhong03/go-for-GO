package handler

import (
	"context"
	"github.com/google/uuid"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/stretchr/testify/mock"
	"github.com/uber/jaeger-client-go"
	"net/http"
	"time"
	"user/dao"
	userProto "user/proto/golang/user"
	"user/tool/jwt"
)

func (u *user) EmailDuplicated(ctx context.Context, req *userProto.EmailDuplicatedRequest, rsp *userProto.EmailDuplicatedResponse) (_ error) {
	if err := u.validate.Struct(req); err != nil {
		rsp.SetStatusAndMsg(http.StatusProxyAuthRequired, err.Error())
		return
	}
	var md metadata.Metadata
	var ok bool
	if md, ok = metadata.FromContext(ctx); !ok || md == nil {
		rsp.SetStatusAndMsg(http.StatusProxyAuthRequired, MessageUnableGetMetadata)
		return
	}
	var xid string
	if xid, ok = md.Get("X-Request-Id"); !ok || xid == "" {
		rsp.SetStatusAndMsg(http.StatusProxyAuthRequired, MessageThereIsNoXReqId)
		return
	}
	if _, err := uuid.Parse(xid); err != nil {
		rsp.SetStatusAndMsg(http.StatusProxyAuthRequired, err.Error())
		return
	}

	// api gateway에서 옳바르지 않은 JWT 필터링 기능 추가 필요 (proto도 변경)
	var id string
	if ss, ok := md.Get("Unique-Authorization"); ok && ss != "" {
		c, err := jwt.ParseDuplicateCertClaimFromJWT(ss)
		if err != nil { rsp.SetStatusAndMsg(http.StatusForbidden, err.Error()); return }
		id = c.UserId
	}

	sps, ok := md.Get("Span-Context")
	if !ok {
		rsp.SetStatusAndMsg(http.StatusProxyAuthRequired, MessageNoSpanContext)
		return
	}
	cs, err := jaeger.ContextFromString(sps)
	if err != nil {
		rsp.SetStatusAndMsg(http.StatusProxyAuthRequired, err.Error())
		return
	}

	var ud dao.UserDAOService
	switch ctx.Value("env") {
	case "test":
		mockStore := ctx.Value("mockStore").(*mock.Mock)
		ud = u.udc.GetTestUserDAO(mockStore)
	default:
		ud = u.udc.GetDefaultUserDAO()
	}

	dsp := u.tracer.StartSpan("CheckIfEmailExist", opentracing.ChildOf(cs))
	dsp.SetTag("X-Request-Id", xid)
	exist, err := ud.CheckIfEmailExist(req.Email)
	dsp.LogFields(log.Bool("exist", exist), log.Error(err))
	dsp.Finish()

	if err != nil {
		rsp.SetStatusAndMsg(http.StatusInternalServerError, err.Error())
		return
	}

	if exist {
		rsp.SetStatusAndMsg(StatusEmailDuplicated, MessageEmailDuplicated)
		return
	}

	ss, err := jwt.GenerateDuplicateCertJWT(id, req.Email, time.Hour)
	if err != nil {
		rsp.SetStatusAndMsg(http.StatusInternalServerError, err.Error())
		return
	}

	rsp.SetStatusAndMsg(http.StatusOK, MessageEmailNotDuplicated)
	rsp.Authorization = ss
	return
}