package mongodb

import (
	"context"
	"github.com/zytekaron/zog-server/src/database"
	"github.com/zytekaron/zog-server/src/types"
	"github.com/zytekaron/zog-server/src/types/find"
	"github.com/zytekaron/zog-server/src/types/updates"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type LogRepository struct {
	logs *mongo.Collection
}

func NewLogRepository(logs *mongo.Collection) database.Controller[*types.Log] {
	return &LogRepository{
		logs: logs,
	}
}

func (l *LogRepository) Insert(ctx context.Context, log *types.Log) error {
	_, err := l.logs.InsertOne(ctx, log)
	if mongo.IsDuplicateKeyError(err) {
		return database.ErrDuplicateKey
	}
	return err
}

func (l *LogRepository) Get(ctx context.Context, id string) (*types.Log, error) {
	var user *types.Log
	err := l.logs.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = database.ErrNoDocuments
		}
		return nil, err
	}

	return user, nil
}

func (l *LogRepository) Update(ctx context.Context, id string, updates updates.Updates[*types.Log]) error {
	res, err := l.logs.UpdateOne(ctx, bson.M{"_id": id}, updates.MongoQuery())
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return database.ErrNoDocuments
	}
	return nil
}

func (l *LogRepository) Delete(ctx context.Context, id string) error {
	res, err := l.logs.DeleteOne(ctx, bson.M{"_id": id})
	if res.DeletedCount == 0 {
		return database.ErrNoDocuments
	}
	return err
}

func (l *LogRepository) Count(ctx context.Context) (int64, error) {
	return l.logs.CountDocuments(ctx, bson.M{})
}

func (l *LogRepository) Find(ctx context.Context, query find.Query[*types.Log], options find.Options[*types.Log]) (database.Iterator[*types.Log], error) {
	cursor, err := l.logs.Find(ctx, query.MongoQuery(), options.MongoOptions())
	if err != nil {
		return nil, err
	}

	iter := &Iterator[*types.Log]{
		out: make(chan *types.Log),
	}
	go func() {
		for cursor.Next(ctx) {
			var log *types.Log
			err := cursor.Decode(&log)
			if err != nil {
				iter.err = err
				break
			}
			iter.out <- log
		}
		close(iter.out)
		cursor.Close(ctx)
	}()

	return iter, nil
}
