// Created by LonelyPale at 2019-12-06

package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type ObjectID = primitive.ObjectID

var NilObjectID = primitive.NilObjectID

func NewObjectID() ObjectID {
	return primitive.NewObjectID()
}

func ObjectIDFromHex(s string) (ObjectID, error) {
	return primitive.ObjectIDFromHex(s)
}
