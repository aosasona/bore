package validation

import (
	"github.com/go-playground/validator/v10"
)

type customValidator struct {
	validateFunc validator.Func
	translation  string
}
