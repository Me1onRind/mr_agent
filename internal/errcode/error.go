package errcode

import (
	"errors"
	"fmt"
)

type Error struct {
	Code int

	message string
	cause   error
}

func NewError(code int, message string) *Error {
	e := &Error{
		Code:    code,
		message: message,
	}
	return e
}

func (e *Error) Error() string {
	return e.message
}

func (e *Error) Unwrap() error {
	return e.cause
}

func (e *Error) Wrap(err error) error {
	return fmt.Errorf("%w, cause:[%w]", e, err)
}

func (e *Error) Withf(msg string, a ...any) error {
	return e.With(fmt.Sprintf(msg, a...))
}

func (e *Error) With(msg string) error {
	return fmt.Errorf("%w, cause:[%s]", e, msg)
}

func ExtractError(err error) *Error {
	var e *Error
	if errors.As(err, &e) {
		return e
	}
	return nil
}
