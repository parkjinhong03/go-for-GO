package subscriber

type changeAuthStatusTest struct {
	AuthId       uint32
	Success      bool
	ExpectError  error
	ExpectMethod map[method]returns
	XRequestID   string
	MessageID	 string
}
