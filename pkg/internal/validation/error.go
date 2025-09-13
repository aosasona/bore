package validation

import "strings"

type validationError struct {
	Field string
	Err   string
}

type ValidationErrors []validationError

func NewValidationError() ValidationErrors {
	return ValidationErrors{}
}

func (v ValidationErrors) Error() string {
	var sb strings.Builder

	for i, ve := range v {
		sb.WriteString("- ")
		sb.WriteString(ve.Field)
		sb.WriteString(": ")
		sb.WriteString(ve.Err)
		if i < len(v)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func (v ValidationErrors) Is(target error) bool {
	_, ok := target.(ValidationErrors)
	return ok
}

func (v *ValidationErrors) Add(field, err string) {
	*v = append(*v, validationError{Field: field, Err: err})
}
