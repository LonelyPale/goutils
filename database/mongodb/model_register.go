package mongodb

import "reflect"

var typeRegistry = make(map[string]reflect.Type)

func RegisterType(typedNil interface{}, names ...string) {
	typed := reflect.TypeOf(typedNil).Elem()
	if len(names) > 0 {
		typeRegistry[names[0]] = typed
	} else {
		typeRegistry[typed.PkgPath()+"."+typed.Name()] = typed
	}
}

func MakeInstance(name string) interface{} {
	if typed, ok := typeRegistry[name]; ok {
		return reflect.New(typed).Elem().Interface()
	}
	return nil
}
