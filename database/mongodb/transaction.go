// Created by LonelyPale at 2019-12-06

package mongodb

import (
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"runtime/debug"
)

type Transaction struct {
	client         *Client
	session        mongo.Session
	sessionContext mongo.SessionContext
	err            error
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

// 自动执行事务，简化过程
func (t *Transaction) Run(fn func(sctx mongo.SessionContext) error) error {
	ctx, cancel := t.client.getContext()
	defer cancel()

	opts := options.Session().SetDefaultReadPreference(readpref.Primary())
	session, err := client.StartSession(opts)
	if err != nil {
		return err
	}
	t.session = session
	defer session.EndSession(ctx)

	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) (e error) {
		t.sessionContext = sessionContext
		defer func() {
			if p := recover(); p != nil {
				str, ok := p.(string)
				if ok {
					e = errors.New(str)
				} else {
					e = errors.New("Transaction: panic recover! ")
				}
				e = t.Rollback()
				debug.PrintStack()
			}
		}()

		e = session.StartTransaction()
		if e != nil {
			return e
		}

		e = fn(sessionContext)
		if e != nil {
			t.err = t.Rollback()
			return e
		} else {
			t.err = t.Commit()
			return t.err
		}
	})

	return err
}
