package dmath

// Cast - casts one numeric type to another. Returned bool indicates if cast was done without overflow.
func Cast[X, Y int | int64 | int32](x X) (Y, bool) {
	var y = Y(x)
	return y, X(y) == x
}
