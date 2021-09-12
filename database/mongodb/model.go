package mongodb

import (
	"context"

	"github.com/lonelypale/goutils/types"
	"github.com/lonelypale/goutils/validator"
)

/*
type Service interface {
	Start() (bool, error)
	OnStart() error

	Stop() bool
	OnStop()

	Reset() (bool, error)
	OnReset() error

	IsRunning() bool

	String() string

	SetLogger(log.Logger)
}
*/
// todo: doc可优化为接口类型,如上例
type Model struct {
	ID         types.ObjectID `bson:"_id,omitempty" json:"id,omitempty" validate:"omitempty,len=12" vCreate:"isdefault" vUpdate:"required" vDelete:"required" vFind:"required" label:"编号"`
	CreateTime types.Time     `bson:"createTime,omitempty" json:"createTime,omitempty" validate:"omitempty" label:"创建时间"`
	UpdateTime types.Time     `bson:"updateTime,omitempty" json:"updateTime,omitempty" validate:"omitempty" label:"更新时间"`

	doc  interface{} `bson:"-" json:"-"`
	coll *Collection `bson:"-" json:"-"`
}

func NewModel(doc interface{}, coll *Collection) Model {
	return Model{doc: doc, coll: coll}
}

//初始化文档和集合
func (m *Model) Init(doc interface{}, coll *Collection) {
	m.doc = doc
	m.coll = coll
}

//是否已初始化
func (m Model) IsInited() bool {
	return m.doc != nil && m.coll != nil
}

func (m Model) Collection() *Collection {
	return m.coll
}

func (m Model) Validate(tags ...string) error {
	return validator.Validate(m.doc, tags...)
}

// Create
func (m Model) Create(ctxs ...context.Context) error {
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

// Update
func (m Model) Update(ctxs ...context.Context) error {
	var ctx context.Context
	if len(ctxs) > 0 {
		ctx = ctxs[0]
	} else {
		var cancel context.CancelFunc
		ctx, cancel = m.coll.GetContext()
		defer cancel()
	}

	if err := m.Validate(UpdateTagName); err != nil {
		return err
	}

	if m.ID.IsZero() {
		return ErrNilObjectID
	}

	filter := ID(m.ID)
	updater := Set(m.doc)
	if _, err := m.coll.UpdateOne(ctx, filter, updater); err != nil {
		return err
	}

	return nil
}

// Delete
func (m Model) Delete(ctxs ...context.Context) error {
	var ctx context.Context
	if len(ctxs) > 0 {
		ctx = ctxs[0]
	} else {
		var cancel context.CancelFunc
		ctx, cancel = m.coll.GetContext()
		defer cancel()
	}

	if m.ID.IsZero() {
		return ErrNilObjectID
	}

	filter := ID(m.ID)
	if _, err := m.coll.DeleteOne(ctx, filter); err != nil {
		return err
	}

	return nil
}

// Find
func (m Model) Find(ctxs ...context.Context) error {
	var ctx context.Context
	if len(ctxs) > 0 {
		ctx = ctxs[0]
	} else {
		var cancel context.CancelFunc
		ctx, cancel = m.coll.GetContext()
		defer cancel()
	}

	if m.ID.IsZero() {
		return ErrNilObjectID
	}

	filter := ID(m.ID)
	if err := m.coll.FindOne(ctx, m.doc, filter); err != nil {
		return err
	}

	return nil
}

func (m Model) Save(ctxs ...context.Context) error {
	var ctx context.Context
	if len(ctxs) > 0 {
		ctx = ctxs[0]
	} else {
		var cancel context.CancelFunc
		ctx, cancel = m.coll.GetContext()
		defer cancel()
	}

	if _, err := m.coll.Save(ctx, m.doc); err != nil {
		return err
	}

	return nil
}

func (m Model) FindOrInsert(filter interface{}, ctxs ...context.Context) error {
	var ctx context.Context
	if len(ctxs) > 0 {
		ctx = ctxs[0]
	} else {
		var cancel context.CancelFunc
		ctx, cancel = m.coll.GetContext()
		defer cancel()
	}

	if _, err := m.coll.FindOrInsert(ctx, filter, m.doc); err != nil {
		return err
	}

	return nil
}
