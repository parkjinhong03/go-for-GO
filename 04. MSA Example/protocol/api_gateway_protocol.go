package protocol

// 읽는 방법: A(파일명)서비스가 B 서비스에게 C를 D하기 위한 Protocol
// Ex) ApiGateway 서비스가 Auth 서비스에게 SignUp을 Request 하기 위한 Protocol
type AuthSignUpRequestProtocol struct {
	Required     RequiredProtocol
	RequestId    string `validate:"required"`
	UserId       string `validate:"required"`
	UserPwd      string `validate:"required"`
	Name         string `validate:"required"`
	PhoneNumber  string `validate:"required"`
	Introduction string
	Email        string `validate:"required"`
}