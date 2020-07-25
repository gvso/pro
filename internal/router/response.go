package router

// ServerErrorResponse represents a 5xx error format.
type ServerErrorResponse struct {
	code    int
	message string
}
