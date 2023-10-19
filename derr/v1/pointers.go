package derr

import "errors"

/*
IsInP - is the same as IsIn, but `err` is a pointer.

`err`:
  - If nil, then result is false by default;
*/
func IsInP(err *error, errs ...error) bool {
	if err == nil {
		err = P(nil)
	}

	return IsIn(*err, errs...)
}

// IsIn - reports whether `err` is one of `errs`.
func IsIn(err error, errs ...error) bool {
	for _, target := range errs {
		if errors.Is(err, target) {
			return true
		}
	}

	return false
}

func P(err error) *error {
	return &err
}
