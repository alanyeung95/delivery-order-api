package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ErrorCode using int to respresent the error
type ErrorCode int

// A list of common expected errors.
const (
	BadRequest       ErrorCode = 400
	ResourceNotFound ErrorCode = 404
	Gone             ErrorCode = 410
	ServerError      ErrorCode = 500
)

// Error specifies the interfaces required by an error in the system.
type Error interface {
	Code() ErrorCode
	Message() string
	error
	json.Marshaler
}

// genericError is an implementation of Error that contains
// an code and error message.
type genericError struct {
	code    ErrorCode
	message string
}

func (e *genericError) Code() ErrorCode {
	return e.code
}

func (e *genericError) Message() string {
	return e.message
}

func (e *genericError) Error() string {
	return fmt.Sprintf("%v : %v", e.code, e.message)
}

func (e *genericError) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Message string `json:"error"`
	}{e.Message()})
}

func (e *genericError) StatusCode() int {
	// http status code reference: https://golang.org/pkg/net/http/
	httpStatus, ok := map[ErrorCode]int{
		ResourceNotFound: http.StatusNotFound,
		BadRequest:       http.StatusBadRequest,
		Gone:             http.StatusGone,
	}[e.Code()]
	if !ok {
		httpStatus = http.StatusInternalServerError
	}
	return httpStatus
}

func NewResourceNotFoundError(err error) Error {
	return &genericError{
		code:    ResourceNotFound,
		message: err.Error(),
	}
}

func NewServerError(err error) Error {
	return &genericError{
		code:    ServerError,
		message: err.Error(),
	}
}

func NewBadRequestError(err error) Error {
	return &genericError{
		code:    BadRequest,
		message: err.Error(),
	}
}

func NewResourceGoneError(err error) Error {
	return &genericError{
		code:    Gone,
		message: err.Error(),
	}
}
