package database

import (
	"context"
	"github.com/zytekaron/zog-server/src/types"
	"github.com/zytekaron/zog-server/src/types/updates"
)

type TokenController interface {
	// Insert inserts a token into the database
	Insert(ctx context.Context, log *types.Token) (err error)

	// Get gets a token from the database
	Get(ctx context.Context, id string) (log *types.Token, err error)

	// Update updates a token in the database
	Update(ctx context.Context, id string, updates *updates.Token) (err error)

	// Delete deletes a token from the database
	Delete(ctx context.Context, id string) (err error)
}
