package entity

type EmailDuplicate struct {
	Email string `json:"email" validate:"required,email,maxLength=30"`
}