package mongodb

import (
	"github.com/LonelyPale/goutils/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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

type Regex = primitive.Regex

type A = types.A //array
type M = types.M //map
type D = types.D //
type E = types.E //
