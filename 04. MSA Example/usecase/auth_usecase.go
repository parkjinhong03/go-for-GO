package usecase

import (
	"MSA.example.com/1/dataservice"
	"MSA.example.com/1/tool/message"
	"github.com/nats-io/nats.go"
)

type authDefaultUseCase struct {
	userD dataservice.UserDataService
	natsM message.NatsMessage
}

func NewAuthDefaultUseCase(userD dataservice.UserDataService, natsM message.NatsMessage) *authDefaultUseCase {
	return &authDefaultUseCase{
		userD: userD,
		natsM: natsM,
	}
}

func (h *authDefaultUseCase) SignUpMsgHandler(msg *nats.Msg) {

}