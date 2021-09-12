package mongodb

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/lonelypale/goutils/types"
)

type CollectionOptions = options.CollectionOptions

type InsertOneOptions = options.InsertOneOptions
type InsertManyOptions = options.InsertManyOptions

type UpdateOptions = options.UpdateOptions

type DeleteOptions = options.DeleteOptions

type FindOneOptions = options.FindOneOptions
type FindOptions = options.FindOptions
type FindOneAndUpdateOptions = options.FindOneAndUpdateOptions

type CountOptions = options.CountOptions

type IndexModel = mongo.IndexModel
type IndexOptions = options.IndexOptions
type CreateIndexesOptions = options.CreateIndexesOptions

type Regex = primitive.Regex

type A = types.A //array
type M = types.M //map
type D = types.D //
type E = types.E //
