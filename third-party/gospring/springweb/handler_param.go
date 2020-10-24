package springweb

import (
	"context"
	"reflect"

	"github.com/go-spring/spring-web"
)

// contextType context.Context 的反射类型
var contextType = reflect.TypeOf((*context.Context)(nil)).Elem()

// webContextType SpringWeb.WebContext 的反射类型
var webContextType = reflect.TypeOf((*SpringWeb.WebContext)(nil)).Elem()

type ParamType uint

const (
	ParamInvalid ParamType = iota
	ParamContext
	ParamWebContext
	ParamJsonStruct
	ParamFormStruct
	ParamUriStruct
	ParamQueryStruct
	ParamHeaderStruct
	ParamOtherStruct
)

type Param struct {
	Type      reflect.Type
	ParamType ParamType
}

func NewParam(typ reflect.Type) Param {
	var paramType ParamType

	switch typ {
	case contextType:
		paramType = ParamContext
	case webContextType:
		paramType = ParamWebContext
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
					paramType = ParamOtherStruct
				}
			} else {
				paramType = ParamOtherStruct
			}
		} else {
			paramType = ParamInvalid
		}
	}

	return Param{
		Type:      typ,
		ParamType: paramType,
	}
}
