package subscriber

import (
	"auth/dao/user"
	"auth/model"
	authProto "auth/proto/golang/auth"
	"auth/tool/random"
	"github.com/google/uuid"
	"time"
)

type changeAuthStatusTest struct {
	AuthId       uint32
	Success      bool
	ExpectError  error
	ExpectMethod map[method]returns
	XRequestID   string
	MessageID	 string
}

func (c changeAuthStatusTest) createTestFromForm(test changeAuthStatusTest) {
	test = c

	if c.AuthId == noneInt 	{ test.AuthId = 0 } 	 else if c.AuthId == 0 		{ test.AuthId = 1 }
	if c.XRequestID == none { test.XRequestID = "" } else if c.XRequestID == "" { test.XRequestID = uuid.New().String() }
	if c.MessageID == none 	{ test.MessageID = "" }  else if c.MessageID == "" 	{ test.MessageID = random.GenerateString(32) }

	if returns, ok := c.ExpectMethod["InsertMessage"]; ok {
		c.setProcessedMessageContext(returns[0].(*model.ProcessedMessage))
	}
	return
}

func (c changeAuthStatusTest) setProcessedMessageContext(psMsg *model.ProcessedMessage) {
	psMsg.ID = psMsgId
	psMsg.MsgId = c.MessageID
	psMsg.CreatedAt = time.Now()
	psMsg.UpdatedAt = time.Now()
	psMsgId++
}

func (c changeAuthStatusTest) setMessageContext(msg *authProto.ChangeAuthStatusMessage) {
	msg.AuthId = c.AuthId
	msg.Success = c.Success
}

func (c changeAuthStatusTest) generateMsgHeader() (header map[string]string) {
	header = make(map[string]string)
	header["XRequestID"] = c.XRequestID
	header["MessageID"] = c.MessageID
	return
}

func (c changeAuthStatusTest) onExpectMethods() {
	for method, returns := range c.ExpectMethod {
		c.onMethod(method, returns)
	}
}

func (c changeAuthStatusTest) onMethod(method method, returns returns) {
	switch method {
	case "InsertMessage":
		mockStore.On("InsertMessage", &model.ProcessedMessage{
			MsgId: c.MessageID,
		}).Return(returns...)
	case "UpdateStatus":
		var status = user.Created
		if !c.Success {
			status = user.Rejected
		}
		mockStore.On("UpdateStatus", c.AuthId, status).Return(returns...)
	case "Commit":
		mockStore.On("Commit").Return(returns...)
	case "Rollback":
		mockStore.On("Rollback").Return(returns...)
	case "Ack":
		mockStore.On("Ack").Return(returns...)
	}
}