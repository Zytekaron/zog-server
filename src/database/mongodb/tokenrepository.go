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

type TokenRepository struct {
	tokens *mongo.Collection
}

func NewTokenRepository(tokens *mongo.Collection) database.Controller[*types.Token] {
	return &TokenRepository{
		tokens: tokens,
	}
}

func (t *TokenRepository) Insert(ctx context.Context, msg *types.Token) error {
	_, err := t.tokens.InsertOne(ctx, msg)
	if mongo.IsDuplicateKeyError(err) {
		return database.ErrDuplicateKey
	}
	return err
}

func (t *TokenRepository) Get(ctx context.Context, id string) (*types.Token, error) {
	var user *types.Token
	err := t.tokens.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = database.ErrNoDocuments
		}
		return nil, err
	}

	return user, nil
}

func (t *TokenRepository) Update(ctx context.Context, id string, updates updates.Updates[*types.Token]) error {
	res, err := t.tokens.UpdateOne(ctx, bson.M{"_id": id}, updates.MongoQuery())
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return database.ErrNoDocuments
	}
	return nil
}

func (t *TokenRepository) Delete(ctx context.Context, id string) error {
	res, err := t.tokens.DeleteOne(ctx, bson.M{"_id": id})
	if res.DeletedCount == 0 {
		return database.ErrNoDocuments
	}
	return err
}

func (t *TokenRepository) Count(ctx context.Context) (int64, error) {
	return t.tokens.CountDocuments(ctx, bson.M{})
}

func (t *TokenRepository) Find(ctx context.Context, query find.Query[*types.Token], options find.Options[*types.Token]) (database.Iterator[*types.Token], error) {
	cursor, err := t.tokens.Find(ctx, query.MongoQuery(), options.MongoOptions())
	if err != nil {
		return nil, err
	}

	iter := &Iterator[*types.Token]{
		out: make(chan *types.Token),
	}
	go func() {
		for cursor.Next(ctx) {
			var log *types.Token
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
