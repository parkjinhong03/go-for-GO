package validator

import (
	"github.com/go-playground/validator/v10"
)

func New() (validate *validator.Validate, err error) {
	validate = validator.New()

	return
}