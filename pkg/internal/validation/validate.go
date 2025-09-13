package validation

import (
	"errors"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	validate   *validator.Validate
	translator ut.Translator
)

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		var name string
		switch {
		case field.Tag.Get("json") != "":
			name, _, _ = strings.Cut(field.Tag.Get("json"), ",")

		case field.Tag.Get("bun") != "":
			name, _, _ = strings.Cut(field.Tag.Get("bun"), ",")
		}

		if name == "-" {
			return ""
		}

		return strings.TrimSpace(name)
	})

	en := en.New()
	uni := ut.New(en, en)

	translator, _ = uni.GetTranslator("en")
	if err := en_translations.RegisterDefaultTranslations(validate, translator); err != nil {
		panic("failed to register default translations: " + err.Error())
	}

	if err := registerCustomValidators(validate, translator); err != nil {
		panic("failed to register custom validators: " + err.Error())
	}
}

func ValidateStruct(s any) error {
	err := validate.Struct(s)

	if err == nil {
		return nil
	}

	switch err := err.(type) {
	case *validator.InvalidValidationError:
		return errors.New("invalid validation error: " + err.Error())
	case validator.ValidationErrors:
		validationErrors := NewValidationError()
		for _, v := range err {
			validationErrors.Add(v.Field(), v.Translate(translator))
		}

		return validationErrors
	default:
		return err
	}
}
