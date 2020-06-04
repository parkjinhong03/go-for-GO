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

const (
	CreatePending string = "CREATE_PENDING"
	Created       string = "CREATED"
	RemovePending string = "REMOVE_PENDING"
	Removed		  string = "REMOVED"
	Reject		  string = "REJECT"
)

type authDefaultUseCase struct {
	userD               dataservice.UserDataService
	validate            *validator.Validate
	apiNatsE, userNatsE natsEncoder.Encoder
}

func NewAuthDefaultUseCase(userD dataservice.UserDataService, validate *validator.Validate,
	apiNatsE, userNatsE natsEncoder.Encoder) *authDefaultUseCase {
		return &authDefaultUseCase{
			userD:     userD,
			validate:  validate,
			apiNatsE:  apiNatsE,
			userNatsE: userNatsE,
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

	p := protocol.ApiGatewaySignUpResponseProtocol{
		Required:   protocol.RequiredProtocol{
			Usage:        "AuthSignUpResponse",
			InputChannel: msg.Reply,
		},
		RequestId:  data.RequestId,
		ResultUser: nil,
		Success:    true,
		ErrorCode:  0,
	}

	_, exist := h.userD.FindByUserId(data.UserId)
	if exist {
		p.Success = false
		p.ErrorCode = UserIdDuplicateErrorCode
		if err := h.apiNatsE.Encode(p); err != nil {
			log.Printf("some error occurs while sending message from auth.signup to api gateway, err: %v\n", err)
		}
		return
	}

	// row 생성 전 동기 호출 응답 추가
	user := model.Users{
		Model:   gorm.Model{},
		UserId:  data.UserId,
		UserPwd: data.UserPwd,
		Status:  CreatePending,
	}
	result, insertErr := h.userD.Insert(&user)
	p.ResultUser = result
	if insertErr != nil {
		log.Printf("unable to insert new user in database, err: %v\n", insertErr)
		errArr := strings.Split(insertErr.Error(), " ")
		errInt, err := strconv.Atoi(errArr[1][:4])
		if errArr[0] != "Error" || err != nil {
			p.ErrorCode = ParsingFailureErrorCode // 에러 코드 파싱 실패
		}
		p.Success = false
		p.ErrorCode = errInt
	}

	if err := h.apiNatsE.Encode(p); err != nil {
		log.Printf("some error occurs while sending message from auth.signup to api gateway, err: %v\n", err)
		h.rejectSignUp(result)
		return
	}
	if insertErr != nil {
		return
	}
	if err := h.userNatsE.Encode(protocol.UserRegistryRequestProtocol{
		Required:     protocol.RequiredProtocol{
			Usage:        "UserRegistryRequest",
			InputChannel: "user.registry",
		},
		RequestId:    data.RequestId,
		ID:           result.ID,
		Name:         data.Name,
		PhoneNumber:  data.PhoneNumber,
		Introduction: data.Introduction,
		Email:        data.Email,
	}); err != nil {
		log.Printf("some error occurs while sending message from auth.signup to user.registry, err:%v\n", err)
		h.rejectSignUp(result)
		return
	}
	return
}

// 사가 트랜잭션 실패했을 경우의 보상 트랜잭션
func (h *authDefaultUseCase) rejectSignUp(user *model.Users) {
	log.Println("executes a compensation transaction because saga transaction has failed.")
	_, _ = h.userD.UpdateStatus(user, Reject)
}

func (h *authDefaultUseCase) RegistryReplyMsgHandler(msg *nats.Msg) {
	data := protocol.AuthRegistryResponseProtocol{}
	if err := json.Unmarshal(msg.Data, &data); err != nil  {
		log.Printf("something occurs error while unmarshal json byte in struct, err: %v\n", err)
		return
	}
	if err := h.validate.Struct(&data); err != nil  {
		log.Printf("something occurs error while validating struct data, err: %v\n", err)
		return
	}
	user, exist := h.userD.Find(data.UserPk)
	if !exist {
		log.Println("There is no user row with that ID.")
		return
	}
	if !data.Success {
		log.Printf("response error from user.registry, error code: %v\n", data.ErrorCode)
		h.rejectSignUp(user)
		return
	}

	if _, err := h.userD.UpdateStatus(user, Created); err != nil {
		h.rejectSignUp(user)
	}
}
