package protocol

import "MSA.example.com/1/model"

type ApiGatewaySignUpResponseProtocol struct {
	Required   RequiredProtocol
	RequestId  string       `validate:"required"`
	ResultUser *model.Users
	Success    bool
	ErrorCode  int
}

type UserRegistryRequestProtocol struct {
	Required     RequiredProtocol
	RequestId    string `validate:"required"`
	ID           uint    `validate:"required"`
	Name         string `validate:"required"`
	PhoneNumber  string `validate:"required"`
	Introduction string
	Email        string `validate:"required"`
}