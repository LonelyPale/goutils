package mongodb

import (
	"context"
	"reflect"

	"go.mongodb.org/mongo-driver/mongo"
)

//分页器接口
//实现接口的方法必须同时是对象 (a A) 或指针 (a *A)，否则断言不到此类型的接口。
type Paginator interface {
	Skip() int64
	Limit() int64
	Result() interface{} //必须是切片指针
	SetTotal(n int64)
}

//判断是否是事务
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
