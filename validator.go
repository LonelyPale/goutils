package goutils

import (
	"reflect"

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
	var vobj reflect.Value

	switch reflect.TypeOf(obj).Kind() {
	case reflect.Ptr:
		vobj = reflect.ValueOf(obj).Elem()
	case reflect.Slice:
		vobj = reflect.ValueOf(obj)
	default:
		return validate(obj, tags...)
	}

	if vobj.Kind() == reflect.Slice {
		length := vobj.Len()
		for i := 0; i < length; i++ {
			o := vobj.Index(i).Interface()
			if err := validate(o, tags...); err != nil {
				return err
			}
		}
	} else {
		return validate(vobj.Interface(), tags...)
	}

	return nil
}

func validate(obj interface{}, tags ...string) error {
	if err := validateStruct(obj); err != nil {
		return err
	}

	for _, tag := range tags {
		if err := validateStruct(obj, tag); err != nil {
			return err
		}
	}

	return nil
}

func validateStruct(obj interface{}, tags ...string) error {
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
