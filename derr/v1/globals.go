package derr

import "errors"

var (
	ErrDone            = errors.New("done")
	ErrExceeded        = errors.New("exceeded")
	ErrDisabled        = errors.New("disabled")
	ErrDuplicated      = errors.New("duplicated")
	ErrNotFound        = errors.New("not found")
	ErrFoundTooMany    = errors.New("found too many")
	ErrUnchanged       = errors.New("unchanged")
	ErrUnauthenticated = errors.New("unauthenticated")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrRolledBack      = errors.New("rolled back")
)
