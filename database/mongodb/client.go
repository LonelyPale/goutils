// Created by LonelyPale at 2019-12-06

package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	opts        *ClientOptions
	mongoClient *mongo.Client
}

func Connect(opts *ClientOptions) (*Client, error) {
	client, err := NewClient(opts)
	if err != nil {
		return nil, err
	}

	ctx, cancel := client.getContext()
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewClient(opts *ClientOptions) (*Client, error) {
	mongoClient, err := mongo.NewClient(opts.mongoClientOptions)
	return &Client{
		opts:        opts,
		mongoClient: mongoClient,
	}, err
}

func (c *Client) MongoClient() *mongo.Client {
	return c.mongoClient
}

func (c *Client) Connect(ctxs ...context.Context) error {
	var ctx context.Context
	if len(ctxs) > 0 && ctxs[0] != nil {
		ctx = ctxs[0]
	} else {
		ctx = context.Background()
	}
	return c.mongoClient.Connect(ctx)
}

func (c *Client) Disconnect(ctxs ...context.Context) error {
	var ctx context.Context
	if len(ctxs) > 0 && ctxs[0] != nil {
		ctx = ctxs[0]
	} else {
		ctx = context.Background()
	}
	return c.mongoClient.Disconnect(ctx)
}

func (c *Client) Database(name string, opts ...*options.DatabaseOptions) *Database {
	return newDatabase(c, name, opts...)
}

func (c *Client) StartSession(opts ...*options.SessionOptions) (mongo.Session, error) {
	return c.mongoClient.StartSession(opts...)
}

func (c *Client) getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), c.opts.Timeout)
}
