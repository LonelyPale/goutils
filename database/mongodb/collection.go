// Created by LonelyPale at 2019-12-06

package mongodb

import (
	"context"
	"errors"
	"github.com/LonelyPale/goutils"
	"github.com/LonelyPale/goutils/database/mongodb/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
)

var ErrResultNil = errors.New("mongodb: the result point cannot be nil")
var ErrResultSlice = errors.New("mongodb: result slice type conversion failure")
var ErrFilterNil = errors.New("mongodb: filter cannot be nil")
var ErrDocumentExists = errors.New("mongodb: document already exists")

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

// 查找一条记录，如果不存在，则插入一条记录
func (coll *Collection) Save(ctx context.Context, filter interface{}, document interface{}, opts ...*options.InsertOneOptions) (types.ObjectID, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = coll.getContext()
		defer cancel()
	}

	if filter == nil {
		return types.NilObjectID, ErrFilterNil
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

// 插入一条记录, 返回 ObjectID
func (coll *Collection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (types.ObjectID, error) {
	mongoCollection := coll.mongoCollection
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = coll.getContext()
		defer cancel()
	}

	res, err := mongoCollection.InsertOne(ctx, document, opts...)
	if err != nil {
		return types.NilObjectID, err
	}

	id := res.InsertedID.(types.ObjectID)
	return id, nil
}

// 插入多条记录，返回 []ObjectID
func (coll *Collection) Insert(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) ([]types.ObjectID, error) {
	mongoCollection := coll.mongoCollection
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = coll.getContext()
		defer cancel()
	}

	res, err := mongoCollection.InsertMany(ctx, documents, opts...)
	if err != nil {
		return nil, err
	}

	var ids []types.ObjectID
	for _, value := range res.InsertedIDs {
		ids = append(ids, value.(types.ObjectID))
	}
	return ids, nil
}

// 修改匹配的第一条记录, 返回修改信息
func (coll *Collection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	mongoCollection := coll.mongoCollection
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = coll.getContext()
		defer cancel()
	}

	res, err := mongoCollection.UpdateOne(ctx, filter, update, opts...)
	if err != nil {
		return nil, err
	}

	//count := int(res.ModifiedCount) //修改记录的数量
	return res, nil
}

// 修改匹配的所有记录, 返回修改信息
func (coll *Collection) Update(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	mongoCollection := coll.mongoCollection
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = coll.getContext()
		defer cancel()
	}

	res, err := mongoCollection.UpdateMany(ctx, filter, update, opts...)
	if err != nil {
		return nil, err
	}

	//count := int(res.ModifiedCount) //修改记录的数量
	return res, nil
}

// 删除匹配的第一条记录, 返回删除数量
func (coll *Collection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (int, error) {
	mongoCollection := coll.mongoCollection
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = coll.getContext()
		defer cancel()
	}

	res, err := mongoCollection.DeleteOne(ctx, filter, opts...)
	if err != nil {
		return 0, err
	}

	count := int(res.DeletedCount)
	return count, nil
}

// 删除所有匹配的记录, 返回删除数量
func (coll *Collection) Delete(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (int, error) {
	mongoCollection := coll.mongoCollection
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = coll.getContext()
		defer cancel()
	}

	res, err := mongoCollection.DeleteMany(ctx, filter, opts...)
	if err != nil {
		return 0, err
	}

	count := int(res.DeletedCount)
	return count, nil
}

// 查找匹配的第一条记录, result 为存储结果的对象指针
func (coll *Collection) FindOne(ctx context.Context, result interface{}, filter interface{}, opts ...*options.FindOneOptions) error {
	mongoCollection := coll.mongoCollection
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = coll.getContext()
		defer cancel()
	}

	if result == nil {
		return ErrResultNil
	}

	err := mongoCollection.FindOne(ctx, filter, opts...).Decode(result)

	return err
}

// 查找匹配的所有记录, result 为存储结果的对象指针
// var users []*model.User 或 users := make([]*model.User, 0) 或 users := []*model.User{}
// 然后 err := Find(&users, collUser, nil)
// users := make([]interface{}, 0) 错误，不能是interface{}万能类型，必须是确定的类型
func (coll *Collection) Find(ctx context.Context, result interface{}, filter interface{}, opts ...*options.FindOptions) error {
	//start := time.Now()

	mongoCollection := coll.mongoCollection
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = coll.getContext()
		defer cancel()
	}

	if result == nil {
		return ErrResultNil
	}

	val := reflect.Indirect(reflect.ValueOf(result))
	//fmt.Printf("val: %v\t%v\n", val, val.Kind())

	if val.Kind() != reflect.Slice {
		return ErrResultSlice
	}

	cur, err := mongoCollection.Find(ctx, filter, opts...)
	defer func() {
		if cur != nil {
			err = cur.Close(ctx)
		}
	}()
	if err != nil {
		return err
	}

	valType := val.Type()
	objType := valType.Elem().Elem()
	//fmt.Printf("typ: %v\t%v\n", typ, objType)

	for cur.Next(ctx) {
		elem := reflect.Indirect(reflect.New(objType)).Addr()
		err := cur.Decode(elem.Interface())
		if err != nil {
			return err
		}
		val = reflect.Append(val, elem)
	}

	if err := cur.Err(); err != nil {
		return err
	}

	//fmt.Printf("val: %v\n", val)
	//fmt.Println(val.Interface().([]*model.User)[0])

	err = goutils.DeepCopy(result, val.Interface()) //深拷贝到源对象

	//log.Println(reflect.TypeOf(result).Elem().Elem().Kind().String())
	//log.Printf("addr of osa:%p,\taddr:%p \t content:%v\n", arr, *arr, *arr)

	//elapsed := time.Since(start)
	//fmt.Println("run elapsed: ", elapsed)

	return err
}

// 统计匹配的记录数量
func (coll *Collection) Count(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int, error) {
	mongoCollection := coll.mongoCollection
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = coll.getContext()
		defer cancel()
	}

	count, err := mongoCollection.CountDocuments(ctx, filter, opts...)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (coll *Collection) getContext() (context.Context, context.CancelFunc) {
	return coll.client.getContext()
}
