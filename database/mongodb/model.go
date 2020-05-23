package mongodb

import (
	"context"
	"reflect"
)

type Model struct {
	doc  interface{}
	coll *Collection
}

func NewModel(doc interface{}, coll *Collection) Model {
	return Model{doc, coll}
}

func (m Model) Collection() *Collection {
	return m.coll
}

func (m Model) Put(ctxs ...context.Context) error {
	var ctx context.Context
	if len(ctxs) > 0 {
		ctx = ctxs[0]
	} else {
		var cancel context.CancelFunc
		ctx, cancel = m.coll.GetContext()
		defer cancel()
	}

	_, err := m.coll.InsertOne(ctx, m.doc)
	if err != nil {
		return err
	}

	return nil
}

func (m Model) Set(ctxs ...context.Context) error {
	var ctx context.Context
	if len(ctxs) > 0 {
		ctx = ctxs[0]
	} else {
		var cancel context.CancelFunc
		ctx, cancel = m.coll.GetContext()
		defer cancel()
	}

	vDoc := reflect.ValueOf(m.doc)
	if vDoc.Kind() == reflect.Ptr {
		vDoc = vDoc.Elem()
	}
	vid := vDoc.FieldByName(IDField)
	if vid.Kind() == reflect.Invalid {
		return ErrNilObjectID
	}

	filter := ID(vid.Interface())
	updater := Set(m.doc)
	if _, err := m.coll.UpdateOne(ctx, filter, updater); err != nil {
		return err
	}

	return nil
}

func (m Model) Get(ctxs ...context.Context) error {
	var ctx context.Context
	if len(ctxs) > 0 {
		ctx = ctxs[0]
	} else {
		var cancel context.CancelFunc
		ctx, cancel = m.coll.GetContext()
		defer cancel()
	}

	vDoc := reflect.ValueOf(m.doc)
	if vDoc.Kind() == reflect.Ptr {
		vDoc = vDoc.Elem()
	}
	vid := vDoc.FieldByName(IDField)
	if vid.Kind() == reflect.Invalid {
		return ErrNilObjectID
	}

	filter := ID(vid.Interface())
	if err := m.coll.FindOne(ctx, m.doc, filter); err != nil {
		return err
	}

	return nil
}

func (m Model) Delete(ctxs ...context.Context) error {
	var ctx context.Context
	if len(ctxs) > 0 {
		ctx = ctxs[0]
	} else {
		var cancel context.CancelFunc
		ctx, cancel = m.coll.GetContext()
		defer cancel()
	}

	vDoc := reflect.ValueOf(m.doc)
	if vDoc.Kind() == reflect.Ptr {
		vDoc = vDoc.Elem()
	}
	vid := vDoc.FieldByName(IDField)
	if vid.Kind() == reflect.Invalid {
		return ErrNilObjectID
	}

	filter := ID(vid.Interface())
	if _, err := m.coll.DeleteOne(ctx, filter); err != nil {
		return err
	}

	return nil
}

func (m Model) Save(filter interface{}, ctxs ...context.Context) error {
	var ctx context.Context
	if len(ctxs) > 0 {
		ctx = ctxs[0]
	} else {
		var cancel context.CancelFunc
		ctx, cancel = m.coll.GetContext()
		defer cancel()
	}

	id, err := m.coll.Save(ctx, filter, m.doc)
	if err != nil {
		return err
	}

	vDoc := reflect.ValueOf(m.doc)
	if vDoc.Kind() == reflect.Ptr {
		vDoc = vDoc.Elem()
	}
	vid := vDoc.FieldByName(IDField)
	if vid.CanSet() {
		for i, n := range id {
			v := vid.Index(i)
			v.SetUint(uint64(n))
		}
	}

	return nil
}
