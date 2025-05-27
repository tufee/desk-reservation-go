package utils

import "fmt"

type BadRequestError struct {
	Message string
}

func (e *BadRequestError) Error() string {
	return e.Message
}

func NewBadRequestError(message string) *BadRequestError {
	return &BadRequestError{Message: message}
}

type InternalServerError struct {
	Message string
	Err     error
}

func (e *InternalServerError) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.Err)
}

func (e *InternalServerError) Unwrap() error {
	return e.Err
}

func NewInternalServerError(message string, err error) *InternalServerError {
	return &InternalServerError{
		Message: message,
		Err:     err,
	}
}
