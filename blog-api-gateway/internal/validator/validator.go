package validator

import "github.com/go-playground/validator/v10"

var (
	validate = validator.New()
)

func Validate(in any) error {
	return validate.Struct(in)
}
