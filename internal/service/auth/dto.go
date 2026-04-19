package auth

import validation "github.com/go-ozzo/ozzo-validation"



type LoginInput struct {
	Email    string
	Password string
}





func (i LoginInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.Email,
			validation.Required,
			validation.Length(3, 100),
		),
		validation.Field(&i.Password,
			validation.Required,
			validation.Length(6, 100),
		),
	)
}