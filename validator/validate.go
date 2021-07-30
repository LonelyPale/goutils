package validator

import (
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"

	"github.com/LonelyPale/goutils/errors"
)

type validate struct {
	validator  *validator.Validate //验证器
	translator ut.Translator       //翻译器
}

func newValidate() *validate {
	return &validate{
		validator: newValidator(),
	}
}

func (v *validate) SetLanguage(language Language) error {
	var err error
	switch language {
	case ZH:
		v.translator, err = SetLanguageZH(v.validator)
	default:
		return errors.New("未找到对应语言的翻译器")
	}
	return err
}

func (v *validate) SetTagName(name string) {
	v.validator.SetTagName(name)
}

func (v *validate) Struct(i interface{}) error {
	if v.translator != nil {
		return Translate(v.validator.Struct(i), v.translator)
	}
	return v.validator.Struct(i)
}

func (v *validate) Var(field interface{}, tag string) error {
	if v.translator != nil {
		return Translate(v.validator.Var(field, tag), v.translator)
	}
	return v.validator.Var(field, tag)
}

func newValidator() *validator.Validate {
	v := validator.New()
	for _, validType := range customValidateTypes {
		v.RegisterCustomTypeFunc(validType.fn, validType.types...)
	}
	return v
}
