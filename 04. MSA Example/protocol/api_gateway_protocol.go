package protocol

import "MSA.example.com/1/model"

type ApiGatewaySignUpResponseProtocol struct {
	Required   RequiredProtocol	`validate:"required"`
	RequestId  string	`validate:"required"`
	ResultUser *model.Users	`validate:"required"`
	Success    bool	`validate:"required"`
	ErrorCode  int	`validate:"required"`
}