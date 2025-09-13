package validation

import (
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go.trulyao.dev/bore/v2/events"
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
	for tag, cv := range validators {
		if err := v.RegisterValidation(tag, cv.validateFunc); err != nil {
			return err
		}

		if strings.TrimSpace(cv.translation) != "" {
			if err := v.RegisterTranslation(
				tag, translator,
				func(ut ut.Translator) error { return ut.Add(tag, cv.translation, true) },
				func(ut ut.Translator, fe validator.FieldError) string {
					t, err := ut.T(tag, fe.Field())
					if err != nil {
						return fe.Error()
					}
					return t
				}); err != nil {
				return err
			}
		}
	}

	return nil
}

func mimetypeValidator(fl validator.FieldLevel) bool {
	_, err := events.MimeTypeFromString(fl.Field().String())
	return err == nil
}
