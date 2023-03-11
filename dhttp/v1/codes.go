package dhttp

/*
CodeError

Format: [N-3][N-2], where
  - 'N-3' - is an HTTP status code;
  - 'N-2' - is an HTTP status sub-code;
*/
type CodeError int

const (
	Code400General         CodeError = 40000
	Code403General         CodeError = 40300
	Code403AccessForbidden CodeError = 40301
	Code404General         CodeError = 40400
	Code500General         CodeError = 50000
)

// IsHTTP - indicates if CodeError belongs to HTTP status code group, where 'code' is valid http status code.
func (c CodeError) IsHTTP(code int) bool {
	return int(c)/100 == code
}
