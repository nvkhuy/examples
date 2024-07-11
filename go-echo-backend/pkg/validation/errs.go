package validation

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// ValidateError error
type ValidateError struct {
	FieldName      string
	ValidateType   string
	ExpectedResult string
	GotResult      interface{}
	Description    string
	fieldError     validator.FieldError
}

// ValidateErrors error
type ValidateErrors struct {
	validator.ValidationErrors
}

func (e *ValidateError) Error() string {
	return e.Description
}

// WithDescription override description
func (e *ValidateError) WithDescription(description string) *ValidateError {
	e.Description = description
	return e
}

// WithTranslate override description
func (e *ValidateError) WithTranslate(translator ut.Translator) *ValidateError {
	return e.WithDescription(e.fieldError.Translate(translator))
}

// TransformError transform error
func (errors ValidateErrors) TransformError() *ValidateError {
	var err = ValidateError{}
	if len(errors.ValidationErrors) > 0 {
		var e = errors.ValidationErrors[0]
		err.FieldName = e.Field()
		err.ValidateType = e.ActualTag()
		err.GotResult = e.Value()
		err.ExpectedResult = e.Param()
		err.fieldError = e

	}

	return &err
}
