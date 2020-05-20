package protocol

import "MSA.example.com/1/model"

type ApiGatewaySignUpResponseProtocol struct {
	Required   RequiredProtocol
	RequestId  string       `validate:"required"`
	ResultUser *model.Users `validate:"required"`
	Success    bool         `validate:"required"`
	ErrorCode  string       `validate:"required"`
}