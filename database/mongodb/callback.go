package mongodb

import (
	"context"
)

type BeforeInserter interface {
	BeforeInsert(ctx context.Context, documents []interface{}) error
}

type AfterInserter interface {
	AfterInsert(ctx context.Context, documents []interface{}) error
}

type BeforeUpdater interface {
	BeforeUpdate(ctx context.Context, filter interface{}, update interface{}) error
}

type AfterUpdater interface {
	AfterUpdate(ctx context.Context, filter interface{}, update interface{}) error
}

type BeforeDeleter interface {
	BeforeDelete(ctx context.Context, filter interface{}) error
}

type AfterDeleter interface {
	AfterDelete(ctx context.Context, filter interface{}) error
}

type BeforeFinder interface {
	BeforeFind(ctx context.Context, result interface{}, filter interface{}) error
}

type AfterFinder interface {
	AfterFind(ctx context.Context, result interface{}, filter interface{}) error
}
