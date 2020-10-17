// Created by LonelyPale at 2019-12-06

package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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

	ctx, cancel := client.GetContext()
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	if err = client.mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return client, nil
}

// NewClient 创建 MongoDB 客户端
func NewClient(opts *ClientOptions) (*Client, error) {
	mongoClient, err := mongo.NewClient(opts.mongoClientOptions)
	if err != nil {
		return nil, err
	}

	return &Client{
		opts:        opts,
		mongoClient: mongoClient,
	}, nil
}

// CloseClient 关闭 MongoDB 客户端
func CloseClient(client *Client) error {
	return client.Disconnect(context.Background())
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

func (c *Client) GetContext() (context.Context, context.CancelFunc) {
	return TimeoutContext(c.opts.Timeout)
}
