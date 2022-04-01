package mongodb

import (
	"context"
	"github.com/zytekaron/zog-server/src/config"
	"github.com/zytekaron/zog-server/src/database"
	"github.com/zytekaron/zog-server/src/types"
	"github.com/zytekaron/zog-server/src/types/updates"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type TokenRepository struct {
	db     *mongo.Database
	tokens *mongo.Collection
}

func NewTokenRepository(db *mongo.Database, dbCfg *config.MongoDB) database.TokenController {
	return &TokenRepository{
		db:     db,
		tokens: db.Collection(dbCfg.TokenCollection),
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

func (t *TokenRepository) Update(ctx context.Context, id string, updates *updates.Token) error {
	res, err := t.tokens.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updates.MongoQuery()})
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
