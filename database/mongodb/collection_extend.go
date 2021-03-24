package mongodb

import (
	"github.com/LonelyPale/goutils/errors"
	"github.com/LonelyPale/goutils/types"
)

func (coll *Collection) DeleteByID(i interface{}) (int64, error) {
	var filter Filter
	switch v := i.(type) {
	case types.ObjectID:
		filter = ID(v)
	case []types.ObjectID:
		filter = In(IDBson, v)
	case string, types.Bytes:
		id, err := types.ObjectIDFrom(v)
		if err != nil {
			return 0, err
		}
		filter = ID(id)
	case []string:
		ids := make([]types.ObjectID, len(v))
		for n, val := range v {
			id, err := types.ObjectIDFrom(val)
			if err != nil {
				return 0, err
			}
			ids[n] = id
		}
		filter = In(IDBson, ids)
	case []types.Bytes:
		ids := make([]types.ObjectID, len(v))
		for n, val := range v {
			id, err := types.ObjectIDFrom(val)
			if err != nil {
				return 0, err
			}
			ids[n] = id
		}
		filter = In(IDBson, ids)
	default:
		return 0, errors.New("invalid argument")
	}

	return coll.Delete(nil, filter)
}

func (coll *Collection) FindByID(result interface{}, i interface{}) error {
	var filter Filter
	switch v := i.(type) {
	case types.ObjectID:
		filter = ID(v)
		return coll.FindOne(nil, result, filter)
	case []types.ObjectID:
		filter = In(IDBson, v)
	case string, types.Bytes:
		id, err := types.ObjectIDFrom(v)
		if err != nil {
			return err
		}
		filter = ID(id)
		return coll.FindOne(nil, result, filter)
	case []string:
		ids := make([]types.ObjectID, len(v))
		for n, val := range v {
			id, err := types.ObjectIDFrom(val)
			if err != nil {
				return err
			}
			ids[n] = id
		}
		filter = In(IDBson, ids)
	case []types.Bytes:
		ids := make([]types.ObjectID, len(v))
		for n, val := range v {
			id, err := types.ObjectIDFrom(val)
			if err != nil {
				return err
			}
			ids[n] = id
		}
		filter = In(IDBson, ids)
	default:
		return errors.New("invalid argument")
	}

	return coll.Find(nil, result, filter)
}
