package customerio

type APIError struct {
	Detail string `json:"detail"`
	Status string `json:"status"`
}

type APIErrors struct {
	Errors []APIError `json:"errors"`
}

func (errs *APIErrors) Error() string {
	for _, e := range errs.Errors {
		return e.Detail
	}

	return "Internal error"
}
