package pdf

type APIError struct {
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return e.Message
}
