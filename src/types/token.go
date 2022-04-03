package types

type Token struct {
	ID        string `json:"id" bson:"_id"`
	OwnerID   string `json:"owner_id" bson:"owner_id"`
	CreatedAt Time   `json:"created_at" bson:"created_at"`
	ExpiresAt Time   `json:"expires_at" bson:"expires_at"`
}
