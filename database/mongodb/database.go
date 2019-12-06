// Created by LonelyPale at 2019-12-06

package mongodb

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	client        *Client
	name          string
	mongoDatabase *mongo.Database
}

func newDatabase(client *Client, name string, opts ...*options.DatabaseOptions) *Database {
	return &Database{
		client:        client,
		name:          name,
		mongoDatabase: client.MongoClient().Database(name, opts...),
	}
}

func (db *Database) MongoDatabase() *mongo.Database {
	return db.mongoDatabase
}

func (db *Database) Client() *Client {
	return db.client
}

func (db *Database) Name() string {
	return db.name
}

func (db *Database) Collection(name string, opts ...*options.CollectionOptions) *Collection {
	return newCollection(db, name, opts...)
}
