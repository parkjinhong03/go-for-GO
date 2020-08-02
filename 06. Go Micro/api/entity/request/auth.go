package request

import (
	authproto "gateway/proto/golang/auth"
)

type UserIdDuplicate struct {
	UserId string `json:"user_id" validate:"required,minLength=4,maxLength=16"`
}

func (ud UserIdDuplicate) ToRequestProto() *authproto.UserIdDuplicatedRequest {
	return &authproto.UserIdDuplicatedRequest{
		UserId: ud.UserId,
	}
}

type UserCreate struct {
	UserId       string `json:"user_id" validate:"required,minLength=4,maxLength=16"`
	UserPw       string `json:"user_pw" validate:"required,minLength=4,maxLength=16"`
	Name         string `json:"name" validate:"required,minLength=2,maxLength=4"`
	PhoneNumber  string `json:"phone_number" validate:"required,strLength=11"`
	Email        string `json:"email" validate:"required,email,maxLength=30"`
	Introduction string `json:"introduction" validate:"maxLength=100"`
}

func (uc UserCreate) ToRequestProto() *authproto.BeforeCreateAuthRequest {
	return &authproto.BeforeCreateAuthRequest{
		UserId:       uc.UserId,
		UserPw:       uc.UserPw,
		Name:         uc.Name,
		PhoneNumber:  uc.PhoneNumber,
		Email:        uc.Email,
		Introduction: uc.Introduction,
	}
}
