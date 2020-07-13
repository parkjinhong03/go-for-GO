package user

import "errors"

var (
	EmailDuplicatedError = errors.New("email duplicated error")
	AuthIdDuplicatedError = errors.New("auth id duplicated error")
	MessageDuplicatedError = errors.New("message duplicated error")
)

var (
	NameTooLongError = errors.New("name too long error")
	EmailTooLongError = errors.New("email too long error")
	PhoneNumberTooLongError = errors.New("phone number too long error")
	MessageTooLongError = errors.New("message id too long error")
)