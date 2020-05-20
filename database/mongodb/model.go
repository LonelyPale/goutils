package mongodb

import (
	"context"
	"errors"
	"reflect"
)

var ErrNilObjectID = errors.New("ObjectID is nil")

type Model struct {
	doc  interface{}
	coll *Collection
}

func NewModel(doc interface{}, coll *Collection) Model {
	return Model{doc, coll}
}

func (m Model) Put(ctxs ...context.Context) error {
	var ctx context.Context
	if len(ctxs) > 0 {
		ctx = ctxs[0]
	}

	id, err := m.coll.InsertOne(ctx, m.doc)
	if err != nil {
		return err
	}

	// 把创建的 ObjectID 写入 doc 中
	vDoc := reflect.ValueOf(m.doc) // 参数必须为指针地址
	eDoc := vDoc.Elem()
	vid := eDoc.FieldByName(IDField)
	if vid.CanSet() {
		for i, n := range id {
			v := vid.Index(i)
			v.SetUint(uint64(n))
		}
	}

	return nil
}

func (m Model) Set(update interface{}, ctxs ...context.Context) error {
	var ctx context.Context
	if len(ctxs) > 0 {
		ctx = ctxs[0]
	}

	vDoc := reflect.ValueOf(m.doc)
	eDoc := vDoc.Elem()
	vid := eDoc.FieldByName(IDField)
	if vid.Kind() == reflect.Invalid {
		return ErrNilObjectID
	}

	filter := NewFilter().ID(vid.Interface())

	if _, err := m.coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (m Model) Get(ctxs ...context.Context) error {
	var ctx context.Context
	if len(ctxs) > 0 {
		ctx = ctxs[0]
	}

	vDoc := reflect.ValueOf(m.doc)
	eDoc := vDoc.Elem()
	vid := eDoc.FieldByName(IDField)
	if vid.Kind() == reflect.Invalid {
		return ErrNilObjectID
	}

	filter := NewFilter().ID(vid.Interface())

	if err := m.coll.FindOne(ctx, m.doc, filter); err != nil {
		return err
	}

	return nil
}

func (m Model) Delete(ctxs ...context.Context) error {
	var ctx context.Context
	if len(ctxs) > 0 {
		ctx = ctxs[0]
	}

	vDoc := reflect.ValueOf(m.doc)
	eDoc := vDoc.Elem()
	vid := eDoc.FieldByName(IDField)
	if vid.Kind() == reflect.Invalid {
		return ErrNilObjectID
	}

	filter := NewFilter().ID(vid.Interface())

	if _, err := m.coll.DeleteOne(ctx, filter); err != nil {
		return err
	}

	return nil
}

func (m Model) Save(filter interface{}, ctxs ...context.Context) error {
	var ctx context.Context
	if len(ctxs) > 0 {
		ctx = ctxs[0]
	}

	id, err := m.coll.Save(ctx, filter, m.doc)
	if err != nil {
		return err
	}

	vDoc := reflect.ValueOf(m.doc)
	eDoc := vDoc.Elem()
	vid := eDoc.FieldByName(IDField)
	if vid.CanSet() {
		for i, n := range id {
			v := vid.Index(i)
			v.SetUint(uint64(n))
		}
	}

	return nil
}
