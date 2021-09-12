package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/lonelypale/goutils/types"
)

type BeforeInserter interface {
	BeforeInsert(ctx context.Context, documents []interface{}, opts interface{}) error
}

type AfterInserter interface {
	AfterInsert(ctx context.Context, documents []interface{}, opts interface{}, ids []types.ObjectID) error
}

type BeforeUpdater interface {
	BeforeUpdate(ctx context.Context, filter interface{}, updater interface{}, opts []*options.UpdateOptions) error
}

type AfterUpdater interface {
	AfterUpdate(ctx context.Context, filter interface{}, updater interface{}, opts []*options.UpdateOptions, result *mongo.UpdateResult) error
}

type BeforeDeleter interface {
	BeforeDelete(ctx context.Context, filter interface{}, opts []*options.DeleteOptions) error
}

type AfterDeleter interface {
	AfterDelete(ctx context.Context, filter interface{}, opts []*options.DeleteOptions, count int64) error
}

type BeforeFinder interface {
	BeforeFind(ctx context.Context, result interface{}, filter interface{}, opts interface{}) error
}

type AfterFinder interface {
	AfterFind(ctx context.Context, result interface{}, filter interface{}, opts interface{}) error
}
