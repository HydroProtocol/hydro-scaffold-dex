package api

import (
	"bytes"
	"fmt"
	"github.com/go-playground/validator"
)

type ApiError struct {
	Code int
	Desc string
}

func NewApiError(code int, desc string) *ApiError {
	return &ApiError{code, desc}
}

func (e *ApiError) Error() string {
	return e.Desc
}

func BindError() *ApiError {
	return &ApiError{Code: -2, Desc: "bind error"}
}

func ValidationError(message string) *ApiError {
	return &ApiError{Code: -1, Desc: message}
}

func MarketNotFoundError(marketID string) *ApiError {
	return &ApiError{Code: -3, Desc: fmt.Sprintf("not support marketID: %s", marketID)}
}

func InvalidPriceAmountError() *ApiError {
	return &ApiError{Code: -4, Desc: fmt.Sprintf("price and amount should be positive number")}
}

func buildErrorMessage(errors validator.ValidationErrors) string {
	buff := bytes.Buffer{}

	for _, err := range errors {
		buff.WriteString(buildSingleError(err))
		buff.WriteString(";")
	}

	return buff.String()
}

func buildSingleError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	default:
		return fmt.Sprintf("Key: '%s' Error:Field validation for '%s' failed on the '%s' tag", err.Namespace(), err.Field(), err.Tag())
	}
}
