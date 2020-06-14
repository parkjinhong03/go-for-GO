package user

import "errors"

var (
	IdDuplicateError = errors.New("user_id duplicate error")
	UnknownError = errors.New("an unknown error")
)