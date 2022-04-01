package types

type Token struct {
	ID        string `json:"id" bson:"_id"`
	OwnerID   string `json:"owner_id" bson:"owner_id"`
	CreatedAt Time   `json:"created_at" bson:"created_at"`
	ExpiresAt Time   `json:"expires_at" bson:"expires_at"`

	Read      bool     `json:"read" bson:"read"`
	ReadLimit int      `json:"read_limit" bson:"read_limit"`
	ReadReset Duration `json:"read_reset" bson:"read_reset"`

	Write      bool     `json:"write" bson:"write"`
	WriteLimit int      `json:"write_limit" bson:"write_limit"`
	WriteReset Duration `json:"write_reset" bson:"write_reset"`
}
