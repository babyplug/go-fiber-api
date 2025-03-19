package validate

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	FailedField string `json:"failedFields"`
	Tag         string `json:"tag"`
	Value       string `json:"value"`
}

func (p *ErrorResponse) String() string {
	return fmt.Sprintf("failedFields=%v,tag=%v,value=%v", p.FailedField, p.Tag, p.Value)
}

func Validate(v any) []*ErrorResponse {
	var errors []*ErrorResponse
	validate := validator.New()
	err := validate.Struct(v)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
