package subscriber

type method string
type returns []interface{}

type createAuthTest struct {
	UserId        string
	UserPw        string
	Name          string
	PhoneNumber   string
	Email         string
	Introduction  string
	XRequestId    string
	ExpectCode    uint32
	ExpectMessage string
	ExpectMethods map[method]returns
}
