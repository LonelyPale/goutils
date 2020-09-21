package mongodb

import (
	"context"
	"reflect"

	"go.mongodb.org/mongo-driver/mongo"
)

func isSessionContext(ctx context.Context) bool {
	if ctx == nil {
		return false
	}

	if _, ok := ctx.(mongo.SessionContext); ok {
		return true
	}

	vctx := reflect.ValueOf(ctx)
	ectx := vctx.Elem()
	if ectx.Kind() == reflect.Struct {
		vContext := ectx.FieldByName("Context")

		switch vContext.Kind() {
		case reflect.Invalid:
			return false
		case reflect.Interface:
			if vContext.CanInterface() {
				iContext := vContext.Interface()
				if c, ok := iContext.(context.Context); ok {
					return isSessionContext(c)
				}
			}
			return false
		default:
			return false
		}
	}

	return false
}
