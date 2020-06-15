package handler

import (
	"auth/dao"
	"auth/dao/user"
	"auth/model"
	"context"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/stretchr/testify/mock"
	"net/http"

	proto "auth/proto/auth"
)

const (
	StatusUserIdDuplicate = 470
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
	switch ctx.Value("env") {
	case "test":
		mockStore := ctx.Value("mockStore").(*mock.Mock)
		ad = e.adc.GetTestAuthDAO(mockStore)
	default:
		ad = e.adc.GetDefaultAuthDAO()
	}

	_, err := ad.Insert(&model.Auth{
		UserId: req.UserId,
		UserPw: req.UserPw,
	})

	switch err {
	case nil:
		rsp.Status = http.StatusCreated
		rsp.Message = "create auth success"
		ad.Commit()
	case user.IdDuplicateError:
		rsp.Status = StatusUserIdDuplicate
		rsp.Message = user.IdDuplicateError.Error()
		ad.Rollback()
	}

	log.Info("Received Auth.CreateAuth request")
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
