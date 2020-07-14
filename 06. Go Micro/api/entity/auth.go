package entity

type UserIdDuplicate struct {
	UserId string `json:"user_id" validate:"required,minLength=4,maxLength=16"`
}