package errs

type SimpleError struct {
	Message string `json:"message"`
}

func NewSimpleError(message string) *SimpleError {
	return &SimpleError{
		Message: message,
	}
}
func (e *SimpleError) Error() string {
	return e.Message
}
