// Created by LonelyPale at 2019-12-06

package types

import (
	"github.com/LonelyPale/goutils"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/LonelyPale/goutils/errors"
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
	i = goutils.PrimitiveValue(i) //去掉指针的包装，以获得原始类型的值

	switch v := i.(type) {
	case ObjectID:
		return v, nil
	case string:
		id, err := ObjectIDFromHex(v)
		if err != nil {
			return NilObjectID, err
		}
		return id, nil
	default:
		return NilObjectID, errors.New("error type")
	}
}
