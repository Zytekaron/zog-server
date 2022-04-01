package find

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Query[T any] interface {
	MongoQuery() bson.M
}

type Options[T any] interface {
	MongoOptions() *options.FindOptions
}
