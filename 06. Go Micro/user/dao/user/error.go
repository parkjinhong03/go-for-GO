package user

import "errors"

var (
	EmailDuplicatedError = errors.New("email duplicated error")
	AuthIdDuplicatedError = errors.New("auth id duplicated error")
	MessageDuplicatedError = errors.New("message duplicated error")
	DataTooLongError = errors.New("data too long error")
)