// Created by LonelyPale at 2019-12-06
package mongodb

import (
	"context"
	"reflect"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/LonelyPale/goutils"
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

	return id, nil
}

func (coll *Collection) Insert(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) ([]types.ObjectID, error) {
	var ids []types.ObjectID

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

	return ids, nil
}

func (coll *Collection) UpdateOne(ctx context.Context, filter interface{}, updater interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	var result *mongo.UpdateResult

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

	return result, nil
}

func (coll *Collection) Update(ctx context.Context, filter interface{}, updater interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	var result *mongo.UpdateResult

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

	return result, nil
}

func (coll *Collection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (int, error) {
	var count int

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

	return count, nil
}

func (coll *Collection) Delete(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (int, error) {
	var count int

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

	return count, nil
}

func (coll *Collection) FindOne(ctx context.Context, result interface{}, filter interface{}, opts ...*options.FindOneOptions) error {
	return coll.findOne(ctx, result, filter, opts...)
}

func (coll *Collection) Find(ctx context.Context, result interface{}, filter interface{}, opts ...*options.FindOptions) error {
	return coll.find(ctx, result, filter, opts...)
}

func (coll *Collection) Count(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int, error) {
	return coll.count(ctx, filter, opts...)
}

// 插入一条记录, 返回 ObjectID
func (coll *Collection) insertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (types.ObjectID, error) {
	if coll == nil {
		return types.NilObjectID, ErrNilCollection
	}

	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = coll.GetContext()
		defer cancel()
	}

	// 根据集合名称找到已注册的模型反射实例
	model := MakeInstance(coll.name)

	// callback before
	//ctx = context.WithValue(ctx, InsertOneOptionsKey, opts)
	before, ok := model.(BeforeInserter)
	if ok {
		if err := before.BeforeInsert(ctx, []interface{}{document}, opts); err != nil {
			return types.NilObjectID, err
		}
	}

	result, err := coll.mongoCollection.InsertOne(ctx, document, opts...)
	if err != nil {
		return types.NilObjectID, err
	}

	id := result.InsertedID.(types.ObjectID)

	// callback after
	//ctx = context.WithValue(ctx, InsertOneResultKey, result)
	after, ok := model.(AfterInserter)
	if ok {
		if err := after.AfterInsert(ctx, []interface{}{document}, opts, []types.ObjectID{id}); err != nil {
			return types.NilObjectID, err
		}
	}

	return id, nil
}

// 插入多条记录，返回 []ObjectID
func (coll *Collection) insert(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) ([]types.ObjectID, error) {
	if coll == nil {
		return nil, ErrNilCollection
	}

	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = coll.GetContext()
		defer cancel()
	}

	model := MakeInstance(coll.name)

	// callback before
	//ctx = context.WithValue(ctx, InsertManyOptionsKey, opts)
	before, ok := model.(BeforeInserter)
	if ok {
		if err := before.BeforeInsert(ctx, documents, opts); err != nil {
			return nil, err
		}
	}

	result, err := coll.mongoCollection.InsertMany(ctx, documents, opts...)
	if err != nil {
		return nil, err
	}

	var ids []types.ObjectID
	for _, value := range result.InsertedIDs {
		ids = append(ids, value.(types.ObjectID))
	}

	// callback after
	//ctx = context.WithValue(ctx, InsertManyResultKey, result)
	after, ok := model.(AfterInserter)
	if ok {
		if err := after.AfterInsert(ctx, documents, opts, ids); err != nil {
			return nil, err
		}
	}

	return ids, nil
}

// 修改匹配的第一条记录, 返回修改信息
func (coll *Collection) updateOne(ctx context.Context, filter interface{}, updater interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	if coll == nil {
		return nil, ErrNilCollection
	}

	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = coll.GetContext()
		defer cancel()
	}

	model := MakeInstance(coll.name)

	// callback before
	//ctx = context.WithValue(ctx, UpdateOptionsKey, opts)
	before, ok := model.(BeforeUpdater)
	if ok {
		if err := before.BeforeUpdate(ctx, filter, updater, opts); err != nil {
			return nil, err
		}
	}

	result, err := coll.mongoCollection.UpdateOne(ctx, filter, updater, opts...)
	if err != nil {
		return nil, err
	}

	// callback after
	//ctx = context.WithValue(ctx, UpdateResultKey, result)
	after, ok := model.(AfterUpdater)
	if ok {
		if err := after.AfterUpdate(ctx, filter, updater, opts, result); err != nil {
			return nil, err
		}
	}

	//count := int(res.ModifiedCount) //修改记录的数量
	return result, nil
}

// 修改匹配的所有记录, 返回修改信息
func (coll *Collection) update(ctx context.Context, filter interface{}, updater interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	if coll == nil {
		return nil, ErrNilCollection
	}

	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = coll.GetContext()
		defer cancel()
	}

	model := MakeInstance(coll.name)

	// callback before
	//ctx = context.WithValue(ctx, UpdateOptionsKey, opts)
	before, ok := model.(BeforeUpdater)
	if ok {
		if err := before.BeforeUpdate(ctx, filter, updater, opts); err != nil {
			return nil, err
		}
	}

	result, err := coll.mongoCollection.UpdateMany(ctx, filter, updater, opts...)
	if err != nil {
		return nil, err
	}

	// callback after
	//ctx = context.WithValue(ctx, UpdateResultKey, result)
	after, ok := model.(AfterUpdater)
	if ok {
		if err := after.AfterUpdate(ctx, filter, updater, opts, result); err != nil {
			return nil, err
		}
	}

	//count := int(res.ModifiedCount) //修改记录的数量
	return result, nil
}

// 删除匹配的第一条记录, 返回删除数量
func (coll *Collection) deleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (int, error) {
	if coll == nil {
		return 0, ErrNilCollection
	}

	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = coll.GetContext()
		defer cancel()
	}

	model := MakeInstance(coll.name)

	// callback before
	//ctx = context.WithValue(ctx, DeleteOptionsKey, opts)
	before, ok := model.(BeforeDeleter)
	if ok {
		if err := before.BeforeDelete(ctx, filter, opts); err != nil {
			return 0, err
		}
	}

	result, err := coll.mongoCollection.DeleteOne(ctx, filter, opts...)
	if err != nil {
		return 0, err
	}

	count := int(result.DeletedCount)

	// callback after
	//ctx = context.WithValue(ctx, DeleteResultKey, result)
	after, ok := model.(AfterDeleter)
	if ok {
		if err := after.AfterDelete(ctx, filter, opts, count); err != nil {
			return 0, err
		}
	}

	return count, nil
}

// 删除所有匹配的记录, 返回删除数量
func (coll *Collection) delete(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (int, error) {
	if coll == nil {
		return 0, ErrNilCollection
	}

	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = coll.GetContext()
		defer cancel()
	}

	model := MakeInstance(coll.name)

	// callback before
	//ctx = context.WithValue(ctx, DeleteOptionsKey, opts)
	before, ok := model.(BeforeDeleter)
	if ok {
		if err := before.BeforeDelete(ctx, filter, opts); err != nil {
			return 0, err
		}
	}

	result, err := coll.mongoCollection.DeleteMany(ctx, filter, opts...)
	if err != nil {
		return 0, err
	}

	count := int(result.DeletedCount)

	// callback after
	//ctx = context.WithValue(ctx, DeleteResultKey, result)
	after, ok := model.(AfterDeleter)
	if ok {
		if err := after.AfterDelete(ctx, filter, opts, count); err != nil {
			return 0, err
		}
	}

	return count, nil
}

// 查找匹配的第一条记录, result 为存储结果的对象指针
func (coll *Collection) findOne(ctx context.Context, result interface{}, filter interface{}, opts ...*options.FindOneOptions) error {
	if coll == nil {
		return ErrNilCollection
	}

	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = coll.GetContext()
		defer cancel()
	}

	if result == nil {
		return ErrNilResult
	}

	model := MakeInstance(coll.name)

	// callback before
	//ctx = context.WithValue(ctx, FindOneOptionsKey, opts)
	before, ok := model.(BeforeFinder)
	if ok {
		if err := before.BeforeFind(ctx, result, filter, opts); err != nil {
			return err
		}
	}

	err := coll.mongoCollection.FindOne(ctx, filter, opts...).Decode(result)
	if err != nil {
		return err
	}

	// callback after
	//ctx = context.WithValue(ctx, SingleResultKey, result)
	after, ok := model.(AfterFinder)
	if ok {
		if err := after.AfterFind(ctx, result, filter, opts); err != nil {
			return err
		}
	}

	return nil
}

// 查找匹配的所有记录, result 为存储结果的对象指针
// var users []*model.User 或 users := make([]*model.User, 0) 或 users := []*model.User{}
// 然后 err := Find(&users, collUser, nil)
// users := make([]interface{}, 0) 错误，不能是interface{}万能类型，必须是确定的类型
func (coll *Collection) find(ctx context.Context, result interface{}, filter interface{}, opts ...*options.FindOptions) error {
	//start := time.Now()
	if coll == nil {
		return ErrNilCollection
	}

	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = coll.GetContext()
		defer cancel()
	}

	if result == nil {
		return ErrNilResult
	}

	model := MakeInstance(coll.name)

	// callback before
	//ctx = context.WithValue(ctx, FindOptionsKey, opts)
	before, ok := model.(BeforeFinder)
	if ok {
		if err := before.BeforeFind(ctx, result, filter, opts); err != nil {
			return err
		}
	}

	val := reflect.Indirect(reflect.ValueOf(result))
	//fmt.Printf("val: %v\t%v\n", val, val.Kind())

	if val.Kind() != reflect.Slice {
		return ErrResultSlice
	}

	cur, err := coll.mongoCollection.Find(ctx, filter, opts...)
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

	//深拷贝到源对象
	if err := goutils.DeepCopy(result, val.Interface()); err != nil {
		return err
	}

	//log.Println(reflect.TypeOf(result).Elem().Elem().Kind().String())
	//log.Printf("addr of osa:%p,\taddr:%p \t content:%v\n", arr, *arr, *arr)

	//elapsed := time.Since(start)
	//fmt.Println("run elapsed: ", elapsed)

	// callback after
	//ctx = context.WithValue(ctx, ResultKey, result)
	after, ok := model.(AfterFinder)
	if ok {
		if err := after.AfterFind(ctx, result, filter, opts); err != nil {
			return err
		}
	}

	return nil
}

// 统计匹配的记录数量
func (coll *Collection) count(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int, error) {
	if coll == nil {
		return 0, ErrNilCollection
	}

	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = coll.GetContext()
		defer cancel()
	}

	count, err := coll.mongoCollection.CountDocuments(ctx, filter, opts...)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

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
