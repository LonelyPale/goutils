package mongodb

import (
	"context"
	"time"
)

const (
	ModelValueKey = "ModelValue"

	InsertOneOptionsKey = "InsertOneOptions"
	InsertOneResultKey  = "InsertOneResult"

	InsertManyOptionsKey = "InsertManyOptions"
	InsertManyResultKey  = "InsertManyResult"

	UpdateOptionsKey = "UpdateOptions"
	UpdateResultKey  = "UpdateResult"

	DeleteOptionsKey = "DeleteOptions"
	DeleteResultKey  = "DeleteResult"

	FindOneOptionsKey = "FindOneOptions"
	SingleResultKey   = "SingleResult"

	FindOptionsKey = "FindOptions"
	ResultKey      = "Result"
)

const DefaultTimeout = 10 * time.Second

func TimeoutContext(ts ...time.Duration) (context.Context, context.CancelFunc) {
	if len(ts) > 0 {
		return context.WithTimeout(context.Background(), ts[0])
	}
	return context.WithTimeout(context.Background(), DefaultTimeout)
}
