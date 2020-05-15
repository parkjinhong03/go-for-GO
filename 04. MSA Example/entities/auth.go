package entities

type AuthSignUpEntities struct {
	UserId string `json:"user_id" validate:"required"`
	UserPwd string `json:"user_pwd" validate:"required"`
	Name string `json:"name" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Introduction string `json:"introduction"`
	Email string `json:"email" validate:"email"`
}