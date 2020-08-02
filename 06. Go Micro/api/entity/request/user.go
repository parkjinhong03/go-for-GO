package request

import (
	userproto "gateway/proto/golang/user"
)

type EmailDuplicate struct {
	Email string `json:"email" validate:"required,email,maxLength=30"`
}

func (ep EmailDuplicate) ToRequestProto() *userproto.EmailDuplicatedRequest {
	return &userproto.EmailDuplicatedRequest{
		Email: ep.Email,
	}
}