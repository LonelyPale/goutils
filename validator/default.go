package validator

import (
	"reflect"

	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"

	"github.com/LonelyPale/goutils/errors"
)

const (
	DefaultTagName      = "validate"
	DefaultLabelTagName = "label"
)

// 用来判断 type T 是否实现了接口 I, 用作类型断言, 如果 T 没有实现接口 I, 则编译错误.
var _ Validator = new(defaultValidator)

// defaultValidator 默认的参数校验器
type defaultValidator struct {
	validator  *validator.Validate //验证器
	translator ut.Translator       //翻译器
}

// NewDefaultValidator defaultValidator 的构造函数
func NewDefaultValidator() *defaultValidator {
	v := &defaultValidator{validator: validator.New()}
	for _, validType := range customValidateTypes {
		v.validator.RegisterCustomTypeFunc(validType.fn, validType.types...)
	}
	return v
}

// Engine 返回原始的参数校验引擎
func (v *defaultValidator) Engine() interface{} {
	return v.validator
}

func (v *defaultValidator) Var(field interface{}, tag string) error {
	return v.validator.Var(field, tag)
}

// Validate 校验参数
func (v *defaultValidator) Validate(obj interface{}, tags ...string) error {
	var vobj reflect.Value

	switch reflect.TypeOf(obj).Kind() {
	case reflect.Ptr:
		vobj = reflect.ValueOf(obj).Elem()
	case reflect.Slice:
		vobj = reflect.ValueOf(obj)
	default:
		return v.validate(obj, tags...)
	}

	if vobj.Kind() == reflect.Slice {
		length := vobj.Len()
		for i := 0; i < length; i++ {
			o := vobj.Index(i).Interface()
			if err := v.validate(o, tags...); err != nil {
				return err
			}
		}
	} else {
		return v.validate(vobj.Interface(), tags...)
	}

	return nil
}

func (v *defaultValidator) validate(obj interface{}, tags ...string) error {
	if err := v.validateStruct(obj); err != nil {
		return err
	}

	for _, tag := range tags {
		if err := v.validateStruct(obj, tag); err != nil {
			return err
		}
	}

	return nil
}

//验证带 tag 的 struct
func (v *defaultValidator) validateStruct(obj interface{}, tags ...string) error {
	if obj == nil {
		return errors.New("validate object is nil")
	}

	//todo: 并发时是否线程安全？
	if len(tags) > 0 && len(tags[0]) > 0 {
		v.validator.SetTagName(tags[0])
	} else {
		v.validator.SetTagName(DefaultTagName)
	}

	if v.translator != nil {
		return Translate(v.validator.Struct(obj), v.translator)
	}

	//err := validate.Struct(u)
	//validationErrors := err.(validator.ValidationErrors)
	return v.validator.Struct(obj)
}

func (v *defaultValidator) SetLanguage(language Language) error {
	var err error
	switch language {
	case ZH:
		v.translator, err = SetLanguageZH(v.validator)
	default:
		return errors.New("未找到对应语言的翻译器")
	}
	return err
}

//自定义验证类型
var customValidateTypes []customValidateType

type customValidateType struct {
	fn    validator.CustomTypeFunc
	types []interface{}
}

func RegisterCustomValidateType(fn validator.CustomTypeFunc, types ...interface{}) {
	customValidateTypes = append(customValidateTypes, customValidateType{fn, types})
}
