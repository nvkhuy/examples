package hubspot

type ApiError struct {
	Status        string `json:"status,omitempty"`
	Message       string `json:"message,omitempty"`
	CorrelationID string `json:"correlationId,omitempty"`
	Category      string `json:"category,omitempty"`
}

func (e *ApiError) Error() string {
	return e.Message
}
