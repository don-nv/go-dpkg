package derr

// Dereference - returns error's pointer value.
func Dereference(err *error) error {
	if err == nil {
		return nil
	}

	return *err
}
