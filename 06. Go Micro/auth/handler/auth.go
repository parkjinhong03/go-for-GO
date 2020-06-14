package handler

import (
	"auth/dao"
	"auth/model"
	"context"
	"github.com/stretchr/testify/mock"
	"net/http"

	log "github.com/micro/go-micro/v2/logger"

	proto "auth/proto/auth"
)

type auth struct{
	adc *dao.AuthDAOCreator
}

func NewAuth(adc *dao.AuthDAOCreator) *auth {
	return &auth{
		adc: adc,
	}
}

func (e *auth) CreateAuth(ctx context.Context, req *proto.CreateAuthRequest, rsp *proto.CreateAuthResponse) error {
	var ad dao.AuthDAOService
	if env := ctx.Value("env"); env == "test" {
		ad = e.adc.GetTestAuthDAO(*ctx.Value("mockStore").(*mock.Mock))
	} else {
		ad = e.adc.GetDefaultAuthDAO()
	}

	if _, err := ad.Insert(&model.Auth{
		UserId: req.UserId,
		UserPw: req.UserPw,
	}); err != nil {
		return err
	}

	ad.Commit()
	log.Info("Received Auth.CreateAuth request")
	rsp.Status = http.StatusOK
	rsp.Message = "test"
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *auth) Stream(ctx context.Context, req *proto.StreamingRequest, stream proto.Auth_StreamStream) error {
	log.Infof("Received Auth.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Infof("Responding: %d", i)
		if err := stream.Send(&proto.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *auth) PingPong(ctx context.Context, stream proto.Auth_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Infof("Got ping %v", req.Stroke)
		if err := stream.Send(&proto.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
