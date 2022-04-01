package mongodb

import (
	"context"
	"github.com/zytekaron/zog-server/src/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient(ctx context.Context, dbCfg *config.MongoDB) (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(dbCfg.URI))
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}

type Iterator[T any] struct {
	out chan T
	err error
}

func (i *Iterator[T]) Chan() chan T {
	return i.out
}

func (i *Iterator[T]) Err() error {
	return i.err
}
