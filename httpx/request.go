package httpx

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type Request struct {
	Vars map[string]string
}

type ValidationMessage struct {
	Field   string
	Message string
}

type ValidationResponse struct {
	Message string
	Errors  []ValidationMessage
}

func (req *Request) Validate(body interface{}) ValidationResponse {

	err := validate.Struct(body)
	var responseBody ValidationResponse

	if err != nil {

		for _, err := range err.(validator.ValidationErrors) {

			msg := ValidationMessage{
				Field:   err.Field(),
				Message: err.Error(),
			}
			responseBody.Errors = append(responseBody.Errors, msg)
		}

	}

	responseBody.Message = "Validation failed"

	return responseBody
}
