package validation

import (
	"regexp"
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go.trulyao.dev/bore/v2/pkg/mimetype"
)

type customValidator struct {
	validateFunc validator.Func
	translation  string
}

var (
	collectionNameRegex = regexp.MustCompile(`^[a-zA-Z0-9 _-]{1,50}$`)

	validators map[string]customValidator = map[string]customValidator{
		"mimetype": {
			validateFunc: mimetypeValidator,
			translation:  "{0} must be a valid MIME type",
		},
		"collection_name": {
			validateFunc: collectionNameValidator,
			translation:  "{0} must be 1-50 characters long and can only contain letters, numbers, spaces, hyphens, and underscores",
		},
	}
)

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
	_, err := mimetype.ParseMimeType(fl.Field().String())
	return err == nil
}

func collectionNameValidator(fl validator.FieldLevel) bool {
	return IsValidCollectionName(fl.Field().String())
}

func IsValidCollectionName(name string) bool {
	name = strings.TrimSpace(name)
	return len(name) > 0 && len(name) <= 50 && collectionNameRegex.MatchString(name)
}
