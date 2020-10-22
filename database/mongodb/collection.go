// Created by LonelyPale at 2019-12-06
package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/LonelyPale/goutils/types"
)

type Collection struct {
	client          *Client
	db              *Database
	name            string
	mongoCollection *mongo.Collection
}

func newCollection(db *Database, name string, opts ...*options.CollectionOptions) *Collection {
	return &Collection{
		client:          db.client,
		db:              db,
		name:            name,
		mongoCollection: db.MongoDatabase().Collection(name, opts...),
	}
}

func (coll *Collection) MongoCollection() *mongo.Collection {
	return coll.mongoCollection
}

func (coll *Collection) Client() *Client {
	return coll.client
}

func (coll *Collection) Database() *Database {
	return coll.db
}

func (coll *Collection) Name() string {
	return coll.name
}

func (coll *Collection) Clone(opts ...*options.CollectionOptions) (*Collection, error) {
	mongoCollection, err := coll.mongoCollection.Clone(opts...)
	if err != nil {
		return nil, err
	}

	return &Collection{
		client:          coll.client,
		db:              coll.db,
		name:            coll.name,
		mongoCollection: mongoCollection,
	}, nil
}

func (coll *Collection) GetContext() (context.Context, context.CancelFunc) {
	return coll.client.GetContext()
}

// 查找一条记录，如果不存在，则插入一条记录
func (coll *Collection) Save(ctx context.Context, filter interface{}, document interface{}, opts ...*options.InsertOneOptions) (types.ObjectID, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = coll.GetContext()
		defer cancel()
	}

	if filter == nil {
		return types.NilObjectID, ErrNilFilter
	}

	count, err := coll.Count(ctx, filter)
	if err != nil {
		return types.NilObjectID, err
	}
	if count > 0 {
		return types.NilObjectID, ErrDocumentExists
	}

	id, err := coll.InsertOne(ctx, document, opts...)
	if err != nil {
		return types.NilObjectID, err
	}

	return id, nil
}

func (coll *Collection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (types.ObjectID, error) {
	var id types.ObjectID

	if coll.client.opts.EnableTransaction {
		if isSessionContext(ctx) {
			return coll.insertOne(ctx, document, opts...)
		} else {
			transaction := NewTransaction(coll.client)
			err := transaction.Run(ctx, func(sctx mongo.SessionContext) error {
				var e error
				id, e = coll.insertOne(sctx, document, opts...)
				return e
			})
			if err != nil {
				return types.NilObjectID, err
			}
		}
	} else {
		return coll.insertOne(ctx, document, opts...)
	}

	return id, nil
}

func (coll *Collection) Insert(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) ([]types.ObjectID, error) {
	var ids []types.ObjectID

	if coll.client.opts.EnableTransaction {
		if isSessionContext(ctx) {
			return coll.insert(ctx, documents, opts...)
		} else {
			transaction := NewTransaction(coll.client)
			err := transaction.Run(ctx, func(sctx mongo.SessionContext) error {
				var e error
				ids, e = coll.insert(sctx, documents, opts...)
				return e
			})
			if err != nil {
				return nil, err
			}
		}
	} else {
		return coll.insert(ctx, documents, opts...)
	}

	return ids, nil
}

func (coll *Collection) UpdateOne(ctx context.Context, filter interface{}, updater interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	var result *mongo.UpdateResult

	if coll.client.opts.EnableTransaction {
		if isSessionContext(ctx) {
			return coll.updateOne(ctx, filter, updater, opts...)
		} else {
			transaction := NewTransaction(coll.client)
			err := transaction.Run(ctx, func(sctx mongo.SessionContext) error {
				var e error
				result, e = coll.updateOne(sctx, filter, updater, opts...)
				return e
			})
			if err != nil {
				return nil, err
			}
		}
	} else {
		return coll.updateOne(ctx, filter, updater, opts...)
	}

	return result, nil
}

func (coll *Collection) Update(ctx context.Context, filter interface{}, updater interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	var result *mongo.UpdateResult

	if coll.client.opts.EnableTransaction {
		if isSessionContext(ctx) {
			return coll.update(ctx, filter, updater, opts...)
		} else {
			transaction := NewTransaction(coll.client)
			err := transaction.Run(ctx, func(sctx mongo.SessionContext) error {
				var e error
				result, e = coll.update(sctx, filter, updater, opts...)
				return e
			})
			if err != nil {
				return nil, err
			}
		}
	} else {
		return coll.update(ctx, filter, updater, opts...)
	}

	return result, nil
}

func (coll *Collection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (int, error) {
	var count int

	if coll.client.opts.EnableTransaction {
		if isSessionContext(ctx) {
			return coll.deleteOne(ctx, filter, opts...)
		} else {
			transaction := NewTransaction(coll.client)
			err := transaction.Run(ctx, func(sctx mongo.SessionContext) error {
				var e error
				count, e = coll.deleteOne(sctx, filter, opts...)
				return e
			})
			if err != nil {
				return 0, err
			}
		}
	} else {
		return coll.deleteOne(ctx, filter, opts...)
	}

	return count, nil
}

func (coll *Collection) Delete(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (int, error) {
	var count int

	if coll.client.opts.EnableTransaction {
		if isSessionContext(ctx) {
			return coll.delete(ctx, filter, opts...)
		} else {
			transaction := NewTransaction(coll.client)
			err := transaction.Run(ctx, func(sctx mongo.SessionContext) error {
				var e error
				count, e = coll.delete(sctx, filter, opts...)
				return e
			})
			if err != nil {
				return 0, err
			}
		}
	} else {
		return coll.delete(ctx, filter, opts...)
	}

	return count, nil
}

func (coll *Collection) FindOne(ctx context.Context, result interface{}, filter interface{}, opts ...*options.FindOneOptions) error {
	return coll.findOne(ctx, result, filter, opts...)
}

//result 必须是指向切片的指针或指向 Pager 接口的指针
func (coll *Collection) Find(ctx context.Context, result interface{}, filter interface{}, opts ...*options.FindOptions) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = coll.GetContext()
		defer cancel()
	}

	pager, ok := result.(Paginator)
	if !ok {
		return coll.find(ctx, result, filter, opts...)
	}

	total, err := coll.Count(ctx, filter)
	if err != nil {
		return err
	}
	pager.SetTotal(total)

	findOptions := options.Find().SetSkip(pager.Skip()).SetLimit(pager.Limit())
	opts = append(opts, findOptions)
	return coll.find(ctx, pager.Result(), filter, opts...)
}

func (coll *Collection) Count(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	return coll.count(ctx, filter, opts...)
}
