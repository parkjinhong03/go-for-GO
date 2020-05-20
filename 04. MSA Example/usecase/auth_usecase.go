package usecase

import (
	"MSA.example.com/1/dataservice"
	"MSA.example.com/1/model"
	"MSA.example.com/1/protocol"
	"MSA.example.com/1/tool/message"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
	"github.com/nats-io/nats.go"
	"log"
)

type authDefaultUseCase struct {
	userD 		dataservice.UserDataService
	natsM 		message.NatsMessage
	validate 	*validator.Validate
}

func NewAuthDefaultUseCase(userD dataservice.UserDataService, natsM message.NatsMessage, validate *validator.Validate) *authDefaultUseCase {
	return &authDefaultUseCase{
		userD: 		userD,
		natsM: 		natsM,
		validate: 	validate,
	}
}

func (h *authDefaultUseCase) SignUpMsgHandler(msg *nats.Msg) {
	data := protocol.AuthSignUpProtocol{}
	if err := json.Unmarshal(msg.Data, &data); err != nil  {
		log.Printf("something occurs error while unmarshal json byte in struct, err: %v\n", err)
		return
	}
	if err := h.validate.Struct(&data); err != nil  {
		log.Printf("something occurs error while validating struct data, err: %v\n", err)
		return
	}

	user := model.Users{
		Model:        gorm.Model{},
		UserId:       data.UserId,
		UserPwd:      data.UserPwd,
	}
	_, err := h.userD.Insert(&user)

	if err != nil {
		log.Printf("unable to insert new user in database, err: %v\n", err)
		_ = h.natsM.Publish(msg.Reply, []byte("test"))
	}
	return
}