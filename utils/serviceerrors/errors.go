package serviceerrors

import (
	"encoding/json"
	"fmt"
)

type ErrorMessage struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details"`
}

var (
	ErrNotFound = ErrorMessage{
		Code:    3,
		Message: "errors.good.notFound",
		Details: struct{}{},
	}

	ErrValidation = ErrorMessage{
		Code:    2,
		Message: "errors.good.incorrect_input",
		Details: struct{}{},
	}

	ErrInternal = ErrorMessage{
		Code:    1,
		Message: "errors.good.internal",
		Details: struct{}{},
	}
)

func (e ErrorMessage) Error() string {
	return fmt.Sprintf("error: code %d, message: %s, details: %v", e.Code, e.Message, e.Details)
}

func (e ErrorMessage) WithDetails(data interface{}) *ErrorMessage {
	e.Details = data
	return &e
}

func (e ErrorMessage) Json() string {
	data, err := json.Marshal(e)

	if err != nil {
		return e.Error()
	}

	return string(data)
}
