package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	STATUS_OK    = "OK"
	STATUS_ERROR = "ERROR"
)

func WriteJson(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

type ErrorResponseStruct struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

func ErrorResponse(err error) ErrorResponseStruct {
	return ErrorResponseStruct{
		Status: STATUS_ERROR,
		Error:  err.Error(),
	}
}

func ValidationError(errs validator.ValidationErrors) ErrorResponseStruct {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is required", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is invalid", err.Field()))
		}
	}
	return ErrorResponseStruct{
		Status: STATUS_ERROR,
		Error:  strings.Join(errMsgs, ", "),
	}
}
