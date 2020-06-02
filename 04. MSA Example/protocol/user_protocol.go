package protocol

import "MSA.example.com/1/model"

type AuthRegistryResponseProtocol struct {
	Required         RequiredProtocol
	RequestId        string `validate:"required"`
	ResultUserInform *model.UserInform
	Success          bool
	ErrorCode        int
}