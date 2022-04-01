package updates

import (
	"go.mongodb.org/mongo-driver/bson"
)

type Updates[T any] interface {
	Apply(T)
	MongoQuery() bson.M
}
