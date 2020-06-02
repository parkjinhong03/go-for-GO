package usecase

import (
	"MSA.example.com/1/dataservice"
	"MSA.example.com/1/model"
	"MSA.example.com/1/protocol"
	natsEncoder "MSA.example.com/1/tool/encoder/nats"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/nats-io/nats.go"
	"log"
	"strconv"
	"strings"
)

type userUseCase struct {
	authNatsE     natsEncoder.Encoder
	userInformDAO dataservice.UserInformDataService
	validate      *validator.Validate
}

func NewUserUseCase(
	authNatsE natsEncoder.Encoder, userInformDAO dataservice.UserInformDataService, validate *validator.Validate) *userUseCase {
		return &userUseCase{
			authNatsE:     authNatsE,
			userInformDAO: userInformDAO,
			validate:      validate,
		}
}

func (u *userUseCase) RegistryMsgHandler(msg *nats.Msg) {
	data := protocol.UserRegistryRequestProtocol{}
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		log.Printf("some error occurs while unmarshal json byte in struct, err: %v\n", err)
		return
	}
	if err := u.validate.Struct(&data); err != nil {
		log.Printf("some error occurs error while validating struct data, err: %v\n", err)
		return
	}

	userInform := model.UserInform{
		UserPk:       data.ID,
		Name:         data.Name,
		PhoneNumber:  data.PhoneNumber,
		Email:        data.Email,
		Introduction: data.Introduction,
	}
	_, err := u.userInformDAO.Insert(&userInform)

	p := protocol.AuthRegistryResponseProtocol{
		Required:  protocol.RequiredProtocol{
			Usage:        "UserRegistryResponse",
			InputChannel: "user.registry.reply",
		},
		RequestId: data.RequestId,
		Success:   true,
		ErrorCode: 0,
	}

	if err != nil {
		p.Success = false
		errArr := strings.Split(err.Error(), " ")
		errInt, err := strconv.Atoi(errArr[1][:4])
		if errArr[0] != "Error" || err != nil {
			p.ErrorCode = ParsingFailureErrorCode // 에러 코드 파싱 실패
		}
		p.ErrorCode = errInt
	}

	if err := u.authNatsE.Encode(p); err != nil {
		// 
	}
}