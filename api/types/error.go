package types

type ErrorResponse struct {
	err error
}

func (e ErrorResponse) Error() error {
	return e.err
}

// Error returns an API response containing a given error
func Error(err error) Response {
	return ErrorResponse{err}
}
