package mongodb

import (
	"context"
	"reflect"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/LonelyPale/goutils"
	"github.com/LonelyPale/goutils/errors"
	"github.com/LonelyPale/goutils/types"
)

// mongodb验证规则: Insert在collection_core中验证，Save在collection中验证，Update在model中验证

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

	if model != nil {
		// 验证新增
		if err := validate(document, CreateTagName); err != nil {
			return types.NilObjectID, err
		}
	}

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

	// 根据集合名称找到已注册的模型反射实例
	model := MakeInstance(coll.name)

	if model != nil {
		// 验证新增
		if err := validate(documents, CreateTagName); err != nil {
			return nil, err
		}
	}

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
func (coll *Collection) deleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (int64, error) {
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

	count := result.DeletedCount

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
func (coll *Collection) delete(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (int64, error) {
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

	count := result.DeletedCount

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

// 查找匹配的所有记录, result 为存储结果的对象指针, 且必须是指针
// var users []model.User 或 users := make([]model.User, 0) 或 users := []model.User{}
// var users []*model.User 或 users := make([]*model.User, 0) 或 users := []*model.User{}
// 然后 err := Find(nil, &users, filter)
// users := make([]interface{}, 0) 错误，不能是interface{}万能类型，必须是确定的类型
func (coll *Collection) find(ctx context.Context, result interface{}, filter interface{}, opts ...*options.FindOptions) (err error) {
	//start := time.Now()

	if coll == nil {
		return ErrNilCollection
	}

	if result == nil {
		return ErrNilResult
	}

	resVal := reflect.ValueOf(result)
	if resVal.Kind() != reflect.Ptr || reflect.Indirect(resVal).Kind() != reflect.Slice {
		return errors.New("result argument must be a pointer to a slice")
	}

	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = coll.GetContext()
		defer cancel()
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

	cur, err := coll.mongoCollection.Find(ctx, filter, opts...)
	defer func() {
		if cur != nil {
			if e := cur.Close(ctx); e != nil {
				err = errors.Errors(err, e)
			}
		}
	}()
	if err != nil {
		return err
	}

	if err := cur.All(ctx, result); err != nil {
		return err
	}

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
func (coll *Collection) count(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
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

	return count, nil
}

// 查找一条记录并修改
func (coll *Collection) findOneAndUpdate(ctx context.Context, result interface{}, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) error {
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

	err := coll.mongoCollection.FindOneAndUpdate(ctx, filter, update, opts...).Decode(result)
	if err != nil {
		return err
	}

	return nil
}

// 查找匹配的所有记录, result 为存储结果的对象指针
// var users []*model.User 或 users := make([]*model.User, 0) 或 users := []*model.User{}
// 然后 err := Find(&users, collUser, nil)
// users := make([]interface{}, 0) 错误，不能是interface{}万能类型，必须是确定的类型
// Deprecated: 弃用，官方已增加 Cursor.All() 方法来设置结果集，不用再自己反射生成结果切片。
func (coll *Collection) findV1(ctx context.Context, result interface{}, filter interface{}, opts ...*options.FindOptions) error {
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
		return errors.New("result slice type conversion failure")
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
