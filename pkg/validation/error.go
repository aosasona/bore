package validation

import "strings"

type validationError struct {
	Field string
	Err   string
}

type ValidationErrors []validationError

var ErrInvalidCollectionName = NewValidationError(
	"collection_name",
	"must be 1-50 characters long and can only contain letters, numbers, spaces, hyphens, and underscores",
)

func NewValidationErrors() ValidationErrors {
	return ValidationErrors{}
}

func NewValidationError(field, err string) ValidationErrors {
	return ValidationErrors{
		{Field: field, Err: err},
	}
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
