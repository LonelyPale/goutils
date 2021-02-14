package mongodb

import (
	"context"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/LonelyPale/goutils/errors"
	"github.com/LonelyPale/goutils/types"
)

// 定义 model 在 coll 内的通用方法，不依赖于 model 的具体值

func (Model) BeforeInsert(ctx context.Context, documents []interface{}, opts interface{}) error {
	for _, doc := range documents {
		vDoc := reflect.ValueOf(doc) // 参数必须为指针地址
		if vDoc.Kind() == reflect.Ptr {
			vDoc = vDoc.Elem()
		} else {
			return errors.ErrMustPointer
		}

		// 设置创建时间
		vCreateTime := vDoc.FieldByName(CreateTimeField)
		if vCreateTime.CanSet() && vCreateTime.Type().String() == "time.Time" {
			now := time.Now()
			vnow := reflect.ValueOf(now)
			vCreateTime.Set(vnow)
		}
	}

	return nil
}

func (Model) AfterInsert(ctx context.Context, documents []interface{}, opts interface{}, ids []types.ObjectID) error {
	for i, doc := range documents {
		vDoc := reflect.ValueOf(doc)
		if vDoc.Kind() == reflect.Ptr {
			vDoc = vDoc.Elem()
		} else {
			return errors.ErrMustPointer
		}

		// 把创建的 ObjectID 写入 doc 中
		vid := vDoc.FieldByName(IDField)
		if vid.CanSet() {
			// SetUint(uint64) 按每字节写入，比直接 Set(Value) 要快。
			for j, n := range ids[i] {
				v := vid.Index(j)
				v.SetUint(uint64(n))
			}
		}
	}

	return nil
}

func (Model) BeforeUpdate(ctx context.Context, filter interface{}, updater interface{}, opts []*options.UpdateOptions) error {
	up, ok := updater.(Updater)
	if !ok {
		return errors.New("not valid Updater")
	}

	set, ok := up[SetKey]
	if !ok {
		return errors.New("not valid Updater[$set]")
	}

	if _, ok := set.(Updater); ok {
		up.UpdateTime(time.Now())
	} else {
		vObj := reflect.ValueOf(set)
		if vObj.Kind() == reflect.Ptr {
			vObj = vObj.Elem()
		} else {
			return errors.ErrMustPointer
		}

		val := vObj.FieldByName(UpdateTimeField)
		if val.CanSet() && val.Type().String() == "time.Time" {
			now := time.Now()
			vnow := reflect.ValueOf(now)
			val.Set(vnow)
		}
	}

	return nil
}
