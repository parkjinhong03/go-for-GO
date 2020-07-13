package user

import "errors"

var (
	BcryptGenerateError = errors.New("bcrypt hash generate error")
	InvalidStatusError = errors.New("this status is invalid")
	NonexistentUserError = errors.New("user does not exist")
)

var (
	UserIdDuplicatedError = errors.New("user id duplicate error")
	MsgIdDuplicateError = errors.New("msg_id duplicate error")
)

var (
	UserIdTooLongError = errors.New("user id too long error")
	MsgIdTooLongError = errors.New("message id too long error")
)