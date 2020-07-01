package user

import "errors"

var (
	EmailDuplicatedError = errors.New("email duplicated error")
	MessageDuplicatedError = errors.New("message duplicated error")
)