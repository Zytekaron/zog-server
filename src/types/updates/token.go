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

	Read      optional.Optional[bool]
	ReadLimit optional.Optional[int]
	ReadReset optional.Optional[types.Duration]

	Write      optional.Optional[bool]
	WriteLimit optional.Optional[int]
	WriteReset optional.Optional[types.Duration]
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

func (t *Token) WithRead(read bool) *Token {
	t.Read = optional.Of(read)
	return t
}

func (t *Token) WithReadLimit(readLimit int) *Token {
	t.ReadLimit = optional.Of(readLimit)
	return t
}

func (t *Token) WithReadReset(readReset types.Duration) *Token {
	t.ReadReset = optional.Of(readReset)
	return t
}

func (t *Token) WithWrite(write bool) *Token {
	t.Write = optional.Of(write)
	return t
}

func (t *Token) WithWriteLimit(writeLimit int) *Token {
	t.WriteLimit = optional.Of(writeLimit)
	return t
}

func (t *Token) WithWriteReset(writeReset types.Duration) *Token {
	t.WriteReset = optional.Of(writeReset)
	return t
}

func (t *Token) Apply(o *types.Token) {
	o.OwnerID = t.OwnerID.OrElse(o.OwnerID)
	o.CreatedAt = t.CreatedAt.OrElse(o.CreatedAt)
	o.ExpiresAt = t.ExpiresAt.OrElse(o.ExpiresAt)

	o.Read = t.Read.OrElse(o.Read)
	o.ReadLimit = t.ReadLimit.OrElse(o.ReadLimit)
	o.ReadReset = t.ReadReset.OrElse(o.ReadReset)

	o.Write = t.Write.OrElse(o.Write)
	o.WriteLimit = t.WriteLimit.OrElse(o.WriteLimit)
	o.WriteReset = t.WriteReset.OrElse(o.WriteReset)
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

	if t.Read.IsPresent() {
		set["read"] = t.Read.Get()
	}
	if t.ReadLimit.IsPresent() {
		set["read_limit"] = t.ReadLimit.Get()
	}
	if t.ReadReset.IsPresent() {
		set["read_reset"] = t.ReadReset.Get()
	}

	if t.Write.IsPresent() {
		set["write"] = t.Write.Get()
	}
	if t.WriteLimit.IsPresent() {
		set["write_limit"] = t.WriteLimit.Get()
	}
	if t.WriteReset.IsPresent() {
		set["write_reset"] = t.WriteReset.Get()
	}

	return bson.M{
		"$set": set,
	}
}
