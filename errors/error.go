package errors

import (
	"errors"
	"fmt"
)

type ErrorWithMessage interface {
	Err() error
	Msg() string
}

type errorWithMessage struct {
	MyErr   error // Wrapped error
	Message string
}

func Is(err error, target error) bool {
	if err == nil {
		return target == nil
	}
	if target == nil {
		return false
	}
	if ewm, ok := err.(ErrorWithMessage); ok {
		return errors.Is(ewm.Err(), target)
	}
	return errors.Is(err, target)
}

func (e *errorWithMessage) Err() error {
	return e.MyErr
}

func (e *errorWithMessage) Msg() string {
	return e.Message
}

func (e *errorWithMessage) Error() string {
	return fmt.Sprintf("%s: %s", e.Message, e.MyErr.Error())
}

func UnWrapperError(err error) ErrorWithMessage {
	wrapper, ok := err.(*errorWithMessage)
	if !ok {
		return &errorWithMessage{
			MyErr:   err,
			Message: "",
		}
	}
	return wrapper
}

func NewWrapperError(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &errorWithMessage{
		Message: msg,
		MyErr:   err,
	}
}
