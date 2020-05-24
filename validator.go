package goutils

import (
	"github.com/go-playground/validator/v10"

	"github.com/LonelyPale/goutils/errors"
)

var customValidateTypes []customValidateType

type customValidateType struct {
	fn    validator.CustomTypeFunc
	types []interface{}
}

func RegisterCustomValidateType(fn validator.CustomTypeFunc, types ...interface{}) {
	customValidateTypes = append(customValidateTypes, customValidateType{fn, types})
}

func Validate(obj interface{}, tags ...string) error {
	if err := validate(obj); err != nil {
		return err
	}

	for _, tag := range tags {
		if err := validate(obj, tag); err != nil {
			return err
		}
	}

	return nil
}

func validate(obj interface{}, tags ...string) error {
	if obj == nil {
		return errors.New("validate object is nil")
	}

	validate := validator.New()
	for _, validType := range customValidateTypes {
		validate.RegisterCustomTypeFunc(validType.fn, validType.types...)
	}

	if len(tags) > 0 && len(tags[0]) > 0 {
		validate.SetTagName(tags[0])
	}

	//err := validate.Struct(u)
	//validationErrors := err.(validator.ValidationErrors)
	return validate.Struct(obj)
}
