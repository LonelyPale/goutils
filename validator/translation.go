package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/zh"
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhtrans "github.com/go-playground/validator/v10/translations/zh"

	"github.com/LonelyPale/goutils/errors"
)

type Language uint

const (
	Invalid Language = iota
	EN
	ZH
)

func SetLanguageZH(v *validator.Validate) (ut.Translator, error) {
	uni := ut.New(zh.New())
	trans, found := uni.GetTranslator("zh")
	if found {
		return nil, errors.New("未找到中文翻译器")
	}

	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		label := field.Tag.Get("label")
		if label == "" {
			return field.Name
		}
		return label
	})

	//验证器注册翻译器
	err := zhtrans.RegisterDefaultTranslations(v, trans)
	if err != nil {
		return nil, err
	}

	return trans, nil
}

func Translate(err error, trans ut.Translator) error {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	errList := make(ValidationErrors, 0)
	for _, e := range errs {
		// can translate each error one at a time.
		fieldName := e.StructField()
		field, _ := e.Type().FieldByName(fieldName)
		jsonTags := strings.Split(field.Tag.Get("json"), ",")
		if len(jsonTags) > 0 && len(jsonTags[0]) > 0 {
			errList = append(errList, NewFieldError(jsonTags[0], e.Translate(trans)))
		} else {
			errList = append(errList, NewFieldError(fieldName, e.Translate(trans)))
		}
	}

	return errList

	//v1: 直接拼接错误字符串
	//var errList []string
	//for _, e := range errs {
	//	// can translate each error one at a time.
	//	errList = append(errList, e.Translate(trans))
	//}
	//
	//return errors.New(strings.Join(errList, "|"))
}
