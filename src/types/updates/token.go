package updates

import (
	"github.com/zytekaron/gotil/v2/optional"
	"github.com/zytekaron/zog-server/src/types"
	"go.mongodb.org/mongo-driver/bson"
)

type Token struct {
	OwnerID   optional.Optional[string]
	CreatedAt optional.Optional[types.Time]
	ExpiresAt optional.Optional[types.Time]
}

func (t *Token) WithLevel(ownerID string) *Token {
	t.OwnerID = optional.Of(ownerID)
	return t
}

func (t *Token) WithCreatedAt(createdAt types.Time) *Token {
	t.CreatedAt = optional.Of(createdAt)
	return t
}

func (t *Token) WithExpiresAt(expiresAt types.Time) *Token {
	t.ExpiresAt = optional.Of(expiresAt)
	return t
}

func (t *Token) Apply(o *types.Token) {
	o.OwnerID = t.OwnerID.OrElse(o.OwnerID)
	o.CreatedAt = t.CreatedAt.OrElse(o.CreatedAt)
	o.ExpiresAt = t.ExpiresAt.OrElse(o.ExpiresAt)
}

func (t *Token) MongoQuery() bson.M {
	set := bson.M{}
	if t.OwnerID.IsPresent() {
		set["owner_id"] = t.OwnerID.Get()
	}
	if t.CreatedAt.IsPresent() {
		set["created_at"] = t.CreatedAt.Get()
	}
	if t.ExpiresAt.IsPresent() {
		set["expires_at"] = t.ExpiresAt.Get()
	}

	return bson.M{
		"$set": set,
	}
}
