package validation

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type customValidator struct {
	validateFunc validator.Func
	translation  string
}

var validators map[string]customValidator = map[string]customValidator{}

func registerCustomValidators(
	v *validator.Validate,
	translator ut.Translator,
) error {
	panic("not implemented")
}
