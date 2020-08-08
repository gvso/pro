package httpresponse

// ErrorResponse represents a 5xx error format.
type ErrorResponse struct {
	code    int
	message string
}

// Error returns a 500 error with the given message.
func Error(msg string) ErrorResponse {
	return ErrorResponse{500, msg}
}
