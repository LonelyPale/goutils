package validator

import "github.com/go-playground/validator/v10"

//自定义验证类型
var customValidateTypes []customValidateType

type customValidateType struct {
	fn    validator.CustomTypeFunc
	types []interface{}
}

func RegisterCustomValidateType(fn validator.CustomTypeFunc, types ...interface{}) {
	customValidateTypes = append(customValidateTypes, customValidateType{fn, types})
}
