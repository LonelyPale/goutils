// Created by LonelyPale at 2019-12-06

package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/lonelypale/goutils/errors"
)

type Transaction struct {
	client         *Client
	session        mongo.Session
	sessionContext mongo.SessionContext
}

func NewTransaction(client *Client) *Transaction {
	return &Transaction{
		client: client,
	}
}

func (t *Transaction) Commit() error {
	return t.session.CommitTransaction(t.sessionContext)
}

func (t *Transaction) Rollback() error {
	return t.session.AbortTransaction(t.sessionContext)
}

// 自动执行事务，简化使用过程。
func (t *Transaction) Run(ctx context.Context, fn func(sctx mongo.SessionContext) error) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = t.client.GetContext()
		defer cancel()
	}

	opts := options.Session().SetDefaultReadPreference(readpref.Primary())
	session, err := t.client.StartSession(opts)
	if err != nil {
		return err
	}

	t.session = session
	defer session.EndSession(ctx)

	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) (e error) {
		t.sessionContext = sessionContext

		defer func() {
			if r := recover(); r != nil || e != nil {
				txerr := t.Rollback()
				e = errors.Errors(e, errors.UnknownError(r), txerr)
				//debug.PrintStack()
			}
		}()

		if e = session.StartTransaction(); e != nil {
			return e
		}

		if e = fn(sessionContext); e != nil {
			return e
		}

		if e = t.Commit(); e != nil {
			return e
		}

		return nil
	})

	return err
}
