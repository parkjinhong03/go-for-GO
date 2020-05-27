package usecase

import (
	"MSA.example.com/1/dataservice"
	"MSA.example.com/1/model"
	"MSA.example.com/1/protocol"
	natsEncoder "MSA.example.com/1/tool/encoder/nats"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
	"github.com/nats-io/nats.go"
	"log"
	"strconv"
	"strings"
)

type authDefaultUseCase struct {
	userD    dataservice.UserDataService
	validate *validator.Validate
	agNatsE  natsEncoder.Encoder
}

func NewAuthDefaultUseCase(
	userD dataservice.UserDataService, validate *validator.Validate, agNatsE natsEncoder.Encoder) *authDefaultUseCase {
	return &authDefaultUseCase{
		userD:    userD,
		validate: validate,
		agNatsE:  agNatsE,
	}
}

func (h *authDefaultUseCase) SignUpMsgHandler(msg *nats.Msg) {
	data := protocol.AuthSignUpRequestProtocol{}
	if err := json.Unmarshal(msg.Data, &data); err != nil  {
		log.Printf("something occurs error while unmarshal json byte in struct, err: %v\n", err)
		return
	}
	if err := h.validate.Struct(&data); err != nil  {
		log.Printf("something occurs error while validating struct data, err: %v\n", err)
		return
	}

	user := model.Users{
		Model:   gorm.Model{},
		UserId:  data.UserId,
		UserPwd: "",
	}
	result, err := h.userD.Insert(&user)

	p := protocol.ApiGatewaySignUpResponseProtocol{
		Required:   protocol.RequiredProtocol{
			Usage:        "AuthSignUpResponse",
			InputChannel: msg.Reply,
		},
		RequestId:  data.RequestId,
		ResultUser: result,
		Success:    true,
		ErrorCode:  0,
	}

	if err != nil {
		log.Printf("unable to insert new user in database, err: %v\n", err)
		errArr := strings.Split(err.Error(), " ")
		errInt, err := strconv.Atoi(errArr[1][:4])
		if errArr[0] != "Error" || err != nil {
			p.ErrorCode = ParsingFailureErrorCode // 에러 코드 파싱 실패
		}
		p.ResultUser = nil
		p.Success = false
		p.ErrorCode = errInt
	}

	err = h.agNatsE.Encode(p)
	if err != nil {
		log.Printf("some error occurs while proccessing that send message from auth.signup, err: %v\n", err)
	}
	return
}