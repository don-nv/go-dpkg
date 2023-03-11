package derr

import (
	"context"
	"errors"
)

// InP - is the same as In, but 'err' is a pointer.
func InP(err *error, errs ...error) bool {
	return In(*err, errs...)
}

// In - reports whether 'err' is one of 'errs'.
func In(err error, errs ...error) bool {
	for _, target := range errs {
		if errors.Is(err, target) {
			return true
		}
	}

	return false
}

/*
CatchF - commonly is used in conjunction with 'defer' to catch an error on return. In this use case 'err' shadowing
has to be avoided.
*/
func CatchF(err *error, f func(err error)) {
	if InP(err, nil) {
		return
	}

	f(*err)
}

// Join - is the same as errors.Join.
func Join(errs ...error) error {
	return errors.Join(errs...)
}

// JoinInP - is the same as Join, but is intended to be used with 'defer' flow control.
func JoinInP(err *error, errs ...error) {
	*err = errors.Join(*err, errors.Join(errs...))
}

// ContextExpired - reports whether 'err' contains one of the standard context expiration errors.
func ContextExpired(err error) bool {
	return In(err, context.Canceled, context.DeadlineExceeded)
}
