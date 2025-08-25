package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	StatusOK    = "OK"
	StatusError = "ERROR"
)

type Response struct {
	Status string `json:"status"` // "Error" or "OK"
	Error  string `json:"error,omitempty"`
	Alias  string `json:"alias,omitempty"`
}

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	errMsgs := make([]string, len(errs), 0)

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field '%s' is required", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("field '%s' is not a valid url", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field '%s' is not valid", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsgs, "; "),
	}
}
