// Created by LonelyPale at 2019-12-06
package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/lonelypale/goutils/errors"
	"github.com/lonelypale/goutils/types"
)

type Collection struct {
	client          *Client
	db              *Database
	name            string
	mongoCollection *mongo.Collection
}

func newCollection(db *Database, name string, opts ...*CollectionOptions) *Collection {
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

func (coll *Collection) Clone(opts ...*CollectionOptions) (*Collection, error) {
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

// todo: 考虑是否用 FindOneAndUpdate 替换，如用 FindOneAndUpdate 替换，则无法调用 Insert、Update callback
// 创建或更新一条记录
func (coll *Collection) Save(ctx context.Context, document interface{}, opts ...*InsertOneOptions) (types.ObjectID, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = coll.GetContext()
		defer cancel()
	}

	var vid reflect.Value
	vDoc := reflect.Indirect(reflect.ValueOf(document))

	switch vDoc.Kind() {
	case reflect.Map:
		vid = vDoc.MapIndex(reflect.ValueOf(IDBson)) //map 没有struct tag(bson)，所以无法从 id 自动转换为 _id
	case reflect.Struct:
		vid = vDoc.FieldByName(IDField)
	}

	//1、没有_id时创建
	if !vid.IsValid() {
		return coll.InsertOne(ctx, document, opts...)
	}

	//2、有_id时修改
	id := vid.Interface().(types.ObjectID)
	filter := ID(id)
	updater := Set(document)

	if err := validate(document, UpdateTagName); err != nil {
		return id, err
	}

	result, err := coll.UpdateOne(ctx, filter, updater)
	if err != nil {
		return id, err
	}

	//2.1、有_id时，但数据库不存在该记录，则创建
	if result.UpsertedID == nil {
		return coll.InsertOne(ctx, document, opts...)
	}

	return id, nil
}

// 推荐优先使用唯一索引来解决字段的唯一性问题，只有在有性能问题时才使用先查询再插入的方法
// 查找一条记录，如果不存在，则插入一条记录
func (coll *Collection) FindOrInsert(ctx context.Context, filter interface{}, document interface{}, opts ...*InsertOneOptions) (types.ObjectID, error) {
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

func (coll *Collection) InsertOne(ctx context.Context, document interface{}, opts ...*InsertOneOptions) (types.ObjectID, error) {
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

func (coll *Collection) Insert(ctx context.Context, documents []interface{}, opts ...*InsertManyOptions) ([]types.ObjectID, error) {
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

func (coll *Collection) UpdateOne(ctx context.Context, filter interface{}, updater interface{}, opts ...*UpdateOptions) (*mongo.UpdateResult, error) {
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

func (coll *Collection) Update(ctx context.Context, filter interface{}, updater interface{}, opts ...*UpdateOptions) (*mongo.UpdateResult, error) {
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

func (coll *Collection) DeleteOne(ctx context.Context, filter interface{}, opts ...*DeleteOptions) (int64, error) {
	var count int64

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

func (coll *Collection) Delete(ctx context.Context, filter interface{}, opts ...*DeleteOptions) (int64, error) {
	var count int64

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

func (coll *Collection) FindOne(ctx context.Context, result interface{}, filter interface{}, opts ...*FindOneOptions) error {
	return coll.findOne(ctx, result, filter, opts...)
}

//result 必须是指向切片的指针或指向 Pager 接口的指针
func (coll *Collection) Find(ctx context.Context, result interface{}, filter interface{}, opts ...*FindOptions) error {
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

func (coll *Collection) Count(ctx context.Context, filter interface{}, opts ...*CountOptions) (int64, error) {
	return coll.count(ctx, filter, opts...)
}

func (coll *Collection) FindOneAndUpdate(ctx context.Context, result interface{}, filter interface{}, update interface{}, opts ...*FindOneAndUpdateOptions) error {
	return coll.findOneAndUpdate(ctx, result, filter, update, opts...)
}

// 创建唯一索引
func (coll *Collection) CreateUniqueIndex(keys ...string) ([]string, error) {
	if len(keys) == 0 {
		return nil, errors.New("CreateUniqueIndex: keys is empty")
	}

	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)

	// 复合索引
	keysDoc := bson.D{}
	for _, key := range keys {
		if strings.HasPrefix(key, "-") { //降序
			//keysDoc = keysDoc.Append(strings.TrimLeft(key, "-"), int32(-1))
			keysDoc = append(keysDoc, bson.E{Key: strings.TrimLeft(key, "-"), Value: int32(-1)})
		} else { //升序
			//keysDoc = keysDoc.Append(key, int32(1))
			keysDoc = append(keysDoc, bson.E{Key: key, Value: int32(1)})
		}
	}

	indexModel := mongo.IndexModel{
		Keys:    keysDoc,
		Options: options.Index().SetUnique(true),
	}

	return coll.CreateIndex(nil, []mongo.IndexModel{indexModel}, opts)
}
