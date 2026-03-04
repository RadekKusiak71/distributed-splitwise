package errors

import (
	"fmt"
	"net/http"
)

type APIError struct {
	StatusCode int `json:"-"`
	Message    any `json:"message"`
}

func (ae *APIError) Error() string {
	return fmt.Sprintf("APIError occured with status %d", ae.StatusCode)
}

func NewAPIError(message any, statusCode int) *APIError {
	return &APIError{StatusCode: statusCode, Message: message}
}

var (
	APIErrInternalServer        = NewAPIError("Internal server error occured while prcessing request.", http.StatusInternalServerError)
	APIErrBadRequest            = NewAPIError("Error occured while prcessing request.", http.StatusBadRequest)
	APIErrUnathorized           = NewAPIError("Unathorized request", http.StatusUnauthorized)
	APIErrMissingIdempotencyKey = NewAPIError("Missing X-Idempotency-Key header", http.StatusBadRequest)
)
