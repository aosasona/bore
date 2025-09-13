package validation

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type customValidator struct {
	validateFunc validator.Func
	translation  string
}

var validators map[string]customValidator = map[string]customValidator{
	"mimetype": {
		validateFunc: mimetypeValidator,
		translation:  "{0} must be a valid MIME type",
	},
}

func registerCustomValidators(
	v *validator.Validate,
	translator ut.Translator,
) error {
	panic("not implemented")
}

func mimetypeValidator(fl validator.FieldLevel) bool {
	panic("not implemented")
}
