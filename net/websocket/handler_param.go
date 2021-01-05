package websocket

import (
	"reflect"
)

var connType = reflect.TypeOf((*Conn)(nil))

var messageType = reflect.TypeOf((*Message)(nil))

type ParamType uint

const (
	ParamInvalid ParamType = iota
	ParamStruct
	ParamOther
	ParamConn
	ParamMessage
)

type Param struct {
	Type      reflect.Type
	ParamType ParamType
}

func NewParam(typ reflect.Type) Param {
	var paramType ParamType

	switch typ {
	case connType:
		paramType = ParamConn
	case messageType:
		paramType = ParamMessage
	default:
		typed := typ
		if typed.Kind() == reflect.Ptr {
			typed = typed.Elem()
		}

		switch typed.Kind() {
		case reflect.Struct:
			paramType = ParamStruct
		default:
			paramType = ParamOther
		}
	}

	return Param{
		Type:      typ,
		ParamType: paramType,
	}
}
