package apigateway

import (
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

// Error represents an API error
type Error struct {
	Code     int    `json:"code"`
	HTTPCode int    `json:"http_code"`
	Message  string `json:"message"`
}

// Error returns the error message
func (e *Error) Error() string {
	return fmt.Sprintf("%s (%d)", e.Message, e.Code)
}

// NewErrorResponse returns an error response
func NewErrorResponse(err error) *events.APIGatewayProxyResponse {
	var knownError *Error
	if errors.As(err, &knownError) {
		return NewJSONResponse(knownError.HTTPCode, knownError)
	}

	return NewJSONResponse(ErrInternalError.HTTPCode, ErrInternalError)
}

// NewInvalidRequestError returns an error for invalid requests that can be rendered in a HTTP response
func NewInvalidRequestError(message string) error {
	return &Error{
		Code:     10400,
		HTTPCode: 400,
		Message:  "Invalid request: " + message,
	}
}

// List of all known errors
var (
	// ErrInternalError is returned when there's an internal error that must be retried
	ErrInternalError = &Error{
		Code:     10500,
		HTTPCode: 500,
		Message:  "Internal server error, try again later",
	}

	// ErrInvalidRequest is returned when the client request is invalid
	ErrInvalidRequest = &Error{
		Code:     10400,
		HTTPCode: 400,
		Message:  "Invalid request",
	}

	// ErrResourceNotFound is returned when the requested resource is not found
	ErrNotFound = &Error{
		Code:     10404,
		HTTPCode: 404,
		Message:  "Resource not found",
	}
)
