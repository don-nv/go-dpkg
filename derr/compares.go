package derr

import "errors"

func Is(source error, targets ...error) bool {
	for _, target := range targets {
		if !errors.Is(source, target) {
			continue
		}

		return true
	}

	return false
}
