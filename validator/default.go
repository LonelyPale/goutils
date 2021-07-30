package validator

import (
	"reflect"

	"github.com/LonelyPale/goutils/errors"
	"github.com/LonelyPale/goutils/ref"
)

const (
	DefaultTagName      = "validate"
	DefaultLabelTagName = "label"
)

// 用来判断 type T 是否实现了接口 I, 用作类型断言, 如果 T 没有实现接口 I, 则编译错误.
var _ Validator = new(defaultValidator)

// defaultValidator 默认的参数校验器
type defaultValidator struct {
	validator *validate      //默认，无tag的验证器
	cache     *validateCache //结构体验证器缓存：非默认tag缓存
	language  Language
}

// NewDefaultValidator defaultValidator 的构造函数
func NewDefaultValidator() *defaultValidator {
	return &defaultValidator{
		validator: newValidate(),
		cache:     newValidateCache(),
	}
}

func (v *defaultValidator) SetLanguage(language Language) error {
	v.language = language
	return v.validator.SetLanguage(language)
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
	if obj == nil {
		return errors.New("validate object is nil")
	}

	vobj := ref.PrimitiveValue(reflect.ValueOf(obj))
	switch vobj.Kind() {
	case reflect.Struct:
		return v.validate(vobj, tags...)
	case reflect.Slice:
		length := vobj.Len()
		for i := 0; i < length; i++ {
			vo := vobj.Index(i)
			if vo.Kind() == reflect.Struct && vo.CanInterface() {
				if err := v.validate(vo, tags...); err != nil {
					return err
				}
			}
		}
	default:
		return errors.Errorf("Validate: unsupported type %s", vobj.Kind().String())
	}

	return nil
}

func (v *defaultValidator) validate(obj reflect.Value, tags ...string) error {
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
func (v *defaultValidator) validateStruct(obj reflect.Value, tags ...string) error {
	var valid *validate

	if len(tags) > 0 && len(tags[0]) > 0 {
		tag := tags[0]
		typ := obj.Type()
		key := typ.PkgPath() + "." + typ.Name() + "-" + tag

		var ok bool
		valid, ok = v.cache.Get(key)
		if !ok {
			valid = v.addValidator(obj, tag)
		}
	} else {
		valid = v.validator
	}

	if valid == nil {
		return errors.New("validator is null")
	}

	return valid.Struct(obj.Interface())
}

func (v *defaultValidator) addValidator(obj reflect.Value, tag string) *validate {
	v.cache.lock.Lock()
	defer v.cache.lock.Unlock()

	typ := obj.Type()
	key := typ.PkgPath() + "." + typ.Name() + "-" + tag

	valid, ok := v.cache.Get(key)
	if ok {
		return valid
	}

	valid = newValidate()
	valid.SetTagName(tag)
	if err := valid.SetLanguage(v.language); err != nil {
		panic(err)
	}

	v.cache.Set(key, valid)
	return valid
}
