package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func parseJSON(r *http.Request, data interface{}) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(data)
}

type ValidationErrorMessage struct {
	Message map[string]string
}

func formValidationErrorMessage(err error) string {
	fieldErrors := err.(validator.ValidationErrors)
	errorObj := ValidationErrorMessage{}
	for _, fieldErr := range fieldErrors {
		errorObj.Message[fieldErr.Field()] = fieldErr.Error()
	}
	res, err := json.Marshal(errorObj)
	if err != nil {
		panic("This should not happen")
	}
	return string(res)
}

type ErrorMessage struct {
	Message string
}

func sendError(w http.ResponseWriter, errorCode int, errorMessage string) {
	res, err := json.Marshal(ErrorMessage{errorMessage})
	if err != nil {
		panic("This should not happen")
	}
	w.WriteHeader(errorCode)
	w.Write(res)
}