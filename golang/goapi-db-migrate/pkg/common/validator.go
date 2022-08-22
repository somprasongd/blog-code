package common

import (
	"errors"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

// use a single instance , it caches struct info
var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
	trans    ut.Translator
)

func init() {
	en := en.New()
	uni = ut.New(en, en)
	// this is usually know or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans, _ = uni.GetTranslator("en")
	validate = validator.New()
	en_translations.RegisterDefaultTranslations(validate, trans)
}

func ValidateDto(s interface{}) error {
	err := validate.Struct(s)
	if err != nil {
		return errors.New(errToMessage(err))
	}
	return nil
}

func errToMessage(err error) (message string) {
	// translate all error at once
	errs := err.(validator.ValidationErrors)
	fields := removeTopStruct(errs.Translate(trans))
	for k, v := range fields {
		message += ", " + k + ": " + v
	}
	return message[2:]
}

func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		// res[strings.ToLower(strings.ReplaceAll(field[strings.Index(field, ".")+1:], ".", "_"))] = err
		res[strings.ToLower(field[strings.Index(field, ".")+1:])] = strings.ToLower(err)
	}
	return res
}
