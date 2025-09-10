package validation

import (
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

	// TODO: register custom validators
}
