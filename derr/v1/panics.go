package derr

import (
	"fmt"
	"runtime/debug"
)

// OnPanic - invokes 'f'. If 'f' panics, invokes 'r' accepting 'err' containing panic value and a call stack.
func OnPanic(f func(), r func(err error)) {
	defer func() {
		v := recover()
		if v != nil {
			r(fmt.Errorf("%+v\n%s", v, debug.Stack()))
		}
	}()

	f()
}

// PanicOnE - panics if 'err' != nil.
func PanicOnE(err error) {
	if err != nil {
		panic(err)
	}
}
