package entities

import "MSA.example.com/1/model"

type SignUpRequestEntities struct {
	UserId string `json:"user_id" validate:"required"`
	UserPwd string `json:"user_pwd" validate:"required"`
	Name string `json:"name" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Introduction string `json:"introduction"`
	Email string `json:"email" validate:"email"`
}

type SignUpResponseEntities struct {
	StatusCode int          `json:"status_code"`
	ResultUser *model.Users `json:"result_user"`
}