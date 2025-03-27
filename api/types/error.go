package types

type ErrorResponse struct {
	msg string
}

func (e ErrorResponse) Error() string {
	return e.msg
}

// Error returns an API response containing a given error
func Error(msg string) Response {
	return ErrorResponse{msg}
}
