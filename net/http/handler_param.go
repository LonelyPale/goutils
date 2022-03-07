package http

import (
	"reflect"

	"github.com/gin-gonic/gin"
)

// ginContextType *gin.Context 的反射类型
var ginContextType = reflect.TypeOf((*gin.Context)(nil)).Elem()

type ParamType uint

const (
	ParamInvalid ParamType = iota
	ParamStruct
	ParamOther
	ParamGinContext
	ParamJsonStruct
	ParamFormStruct
	ParamUriStruct
	ParamQueryStruct
	ParamHeaderStruct
)

type Param struct {
	Type      reflect.Type
	ParamType ParamType
}

func NewParam(typ reflect.Type) Param {
	var paramType ParamType

	switch typ {
	case ginContextType:
		paramType = ParamGinContext
	default:
		if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct {
			ts := typ.Elem()
			if ts.NumField() > 0 {
				tag := ts.Field(0).Tag
				if s := tag.Get("json"); len(s) > 0 {
					paramType = ParamJsonStruct
				} else if s := tag.Get("form"); len(s) > 0 {
					paramType = ParamFormStruct
				} else if s := tag.Get("uri"); len(s) > 0 {
					paramType = ParamUriStruct
				} else if s := tag.Get("query"); len(s) > 0 {
					paramType = ParamQueryStruct
				} else if s := tag.Get("header"); len(s) > 0 {
					paramType = ParamHeaderStruct
				} else {
					paramType = ParamStruct
				}
			} else {
				paramType = ParamStruct
			}
		} else {
			paramType = ParamOther
		}
	}

	return Param{
		Type:      typ,
		ParamType: paramType,
	}
}
