package protocol

type AuthRegistryResponseProtocol struct {
	Required  RequiredProtocol
	RequestId string `validate:"required"`
	Success   bool
	ErrorCode int
}