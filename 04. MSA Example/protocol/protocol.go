package protocol

type RequiredProtocol struct {
	Usage        string `validate:"required"`
	InputChannel string `validate:"required"`
}