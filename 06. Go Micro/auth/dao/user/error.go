package user

import "errors"

var (
	BcryptGenerateError = errors.New("bcrypt hash generate error")
	IdDuplicateError = errors.New("user_id duplicate error")
	MsgIdDuplicateError = errors.New("msg_id duplicate error")
	DataLengthOverError = errors.New("data too long for column error")
	UnknownError = errors.New("an unknown error")
)