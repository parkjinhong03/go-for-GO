package subscriber

import (
	"auth/dao"
	"auth/dao/user"
	"auth/model"
	proto "auth/proto/auth"
	"auth/tool/random"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/micro/go-micro/v2/broker"
	userProto "user/proto/user"
	"user/subscriber"
)

func (m *auth) CreateAuth(e broker.Event) error {
	header := e.Message().Header
	msgId := header["MessageID"]
	xReqId := header["XRequestID"]

	// msgId 길이 최대 32 제한 추가
	if xReqId == "" || msgId == "" {
		return ErrorForbidden
	}

	if _, err := uuid.Parse(xReqId); err != nil {
		return ErrorForbidden
	}

	recvMsg := proto.CreateAuthMessage{}
	if err := json.Unmarshal(e.Message().Body, &recvMsg); err != nil {
		// 에러 기록
		return ErrorBadRequest
	}

	if err := m.validate.Struct(&recvMsg); err != nil {
		// 에러 기록
		return ErrorBadRequest
	}

	var ad dao.AuthDAOService
	var aftMsgId string
	env, ok := header["Env"]
	if ok && env == "Test" {
		ad = m.adc.GetTestAuthDAO(e.(*CustomEvent).mock)
		aftMsgId = header["AfterMessageID"]
	} else {
		ad = m.adc.GetDefaultAuthDAO()
		aftMsgId = random.GenerateString(32)
	}

	if _, err := ad.InsertMessage(&model.ProcessedMessage{
		MsgId: msgId,
	}); err != nil {
		if err == user.MsgIdDuplicateError { return ErrorDuplicated }
		return nil
	}

	var result *model.Auth
	var err error
	if result, err = ad.InsertAuth(&model.Auth{
		UserId: recvMsg.UserId,
		UserPw: recvMsg.UserPw,
		Status: user.CreatePending,
	}); err != nil {
		ad.Rollback()
		// err 이유 떄문에 insert auth 트랜잭션을 롤백시켰다는 기록
		return nil
	}
	ad.Commit()

	header = make(map[string]string)
	header["XRequestID"] = xReqId
	header["MessageID"] = aftMsgId

	sendMsg := userProto.CreateUserMessage{
		Id:           uint32(result.ID),
		Name:         recvMsg.Name,
		PhoneNumber:  recvMsg.PhoneNumber,
		Email:        recvMsg.Email,
		Introduction: recvMsg.Introduction,
	}

	body, err := json.Marshal(sendMsg)
	if err != nil {
		// 거의 발생 X, 트랜잭션까지는 했지만 publish 하기 전에 메시지 인코딩에서 에러가 발생했다는 기록
		return nil
	}

	if err := m.mq.Publish(subscriber.CreateUserEventTopic, &broker.Message{
		Header: header,
		Body:   body,
	}); err != nil {
		// 트랜잭션까지는 했지만 publish를 못했다는 기록
		return nil
	}

	if err := e.Ack(); err != nil {
		// publish 까지는 했지만 ack는 못했다는 기록
		return nil
	}

	// 정상 처리 기록
	return nil
}