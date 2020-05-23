package mongodb

import (
	"context"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/LonelyPale/goutils/database/mongodb/types"
)

func (m Model) BeforeInsert(ctx context.Context, documents []interface{}, opts interface{}) error {
	for _, doc := range documents {
		vDoc := reflect.ValueOf(doc) // 参数必须为指针地址
		if vDoc.Kind() != reflect.Ptr {
			return ErrMustPointer
		}

		// 设置创建时间
		eDoc := vDoc.Elem()
		vCreateTime := eDoc.FieldByName(CreateTimeField)
		if vCreateTime.CanSet() && vCreateTime.Type().String() == "time.Time" {
			now := time.Now()
			vnow := reflect.ValueOf(now)
			vCreateTime.Set(vnow)
		}
	}

	return nil
}

func (m Model) AfterInsert(ctx context.Context, documents []interface{}, opts interface{}, ids []types.ObjectID) error {
	for i, doc := range documents {
		vDoc := reflect.ValueOf(doc)
		if vDoc.Kind() != reflect.Ptr {
			return ErrMustPointer
		}

		// 把创建的 ObjectID 写入 doc 中
		eDoc := vDoc.Elem()
		vid := eDoc.FieldByName(IDField)
		if vid.CanSet() {
			for j, n := range ids[i] {
				v := vid.Index(j)
				v.SetUint(uint64(n))
			}
		}
	}

	return nil
}

func (m Model) BeforeUpdate(ctx context.Context, filter interface{}, updater interface{}, opts []*options.UpdateOptions) error {
	return nil
}
