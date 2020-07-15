package entity

type UserIdDuplicate struct {
	UserId string `json:"user_id" validate:"required,minLength=4,maxLength=16"`
}

type UserCreate struct {
	UserId       string `json:"user_id" validate:"required,minLength=4,maxLength=16"`
	UserPw       string `json:"user_pw" validate:"required,minLength=4,maxLength=16"`
	Name         string `json:"name" validate:"required,minLength=2,maxLength=4"`
	PhoneNumber  string `json:"phone_number" validate:"required,strLength=11"`
	Email        string `json:"email" validate:"required,email,maxLength=30"`
	Introduction string `json:"introduction" validate:"maxLength=100"`
}