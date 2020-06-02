package protocol

type UserRegistryPublishProtocol struct {
	Required     RequiredProtocol
	RequestId    string `validate:"required"`
	ID           uint    `validate:"required"`
	Name         string `validate:"required"`
	PhoneNumber  string `validate:"required"`
	Introduction string
	Email        string `validate:"required"`
}