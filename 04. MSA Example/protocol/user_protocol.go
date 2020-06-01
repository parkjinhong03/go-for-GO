package protocol

type UserRegistryPublishProtocol struct {
	Required     RequiredProtocol
	RequestId    string `validate:"required"`
	Id           int    `validate:"required"`
	Name         string `validate:"required"`
	PhoneNumber  string `validate:"required"`
	Introduction string
	Email        string `validate:"required"`
}