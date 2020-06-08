package handler

import (
	"context"

	log "github.com/micro/go-micro/v2/logger"

	auth "auth/proto/auth"
)

type Auth struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Auth) Call(ctx context.Context, req *auth.Request, rsp *auth.Response) error {
	log.Info("Received Auth.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Auth) Stream(ctx context.Context, req *auth.StreamingRequest, stream auth.Auth_StreamStream) error {
	log.Infof("Received Auth.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Infof("Responding: %d", i)
		if err := stream.Send(&auth.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Auth) PingPong(ctx context.Context, stream auth.Auth_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Infof("Got ping %v", req.Stroke)
		if err := stream.Send(&auth.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
