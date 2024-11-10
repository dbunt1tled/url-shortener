package http

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

type Response struct {
	Status string `json:"status"`
	Data   any    `json:"data"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOk    = "ok"
	StatusError = "error"
)

func Ok(data any) Response {
	return Response{
		Status: StatusOk,
		Data:   data,
	}
}

func Error(error string) Response {
	return Response{
		Status: StatusError,
		Error:  error,
	}
}

func ValidationErrorResponse(error validator.ValidationErrors) Response {
	var errors []string
	for _, err := range error {
		errors = append(errors, fmt.Sprintf("%s: %s", err.Field(), err.ActualTag()))
	}
	return Response{
		Status: StatusError,
		Error:  strings.Join(errors, ","),
	}

}
