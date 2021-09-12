// Created by LonelyPale at 2019-12-06

package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/lonelypale/goutils/errors"
)

type ObjectID = primitive.ObjectID

var NilObjectID = primitive.NilObjectID

func NewObjectID() ObjectID {
	return primitive.NewObjectID()
}

func ObjectIDFromHex(s string) (ObjectID, error) {
	return primitive.ObjectIDFromHex(s)
}

func ObjectIDFrom(i interface{}) (ObjectID, error) {
	switch v := i.(type) {
	case ObjectID:
		return v, nil
	case *ObjectID:
		return *v, nil
	case string:
		return primitive.ObjectIDFromHex(v)
	case *string:
		return primitive.ObjectIDFromHex(*v)
	default:
		return NilObjectID, errors.New("error type")
	}
}

func ObjectIDSliceFrom(i interface{}) ([]ObjectID, error) {
	switch v := i.(type) {
	case []ObjectID:
		return v, nil
	case []string:
		ids := make([]ObjectID, len(v))
		for n, val := range v {
			id, err := primitive.ObjectIDFromHex(val)
			if err != nil {
				return nil, err
			}
			ids[n] = id
		}
		return ids, nil
	default:
		return nil, errors.New("error type")
	}
}
