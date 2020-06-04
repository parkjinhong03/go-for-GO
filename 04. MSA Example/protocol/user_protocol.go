package protocol

type AuthRegistryResponseProtocol struct {
	Required  RequiredProtocol
	RequestId string `validate:"required"`
	UserPk    uint
	Success   bool
	ErrorCode int
}