package database

import (
	"context"
	"github.com/zytekaron/zog-server/src/types"
	"github.com/zytekaron/zog-server/src/types/find"
	"github.com/zytekaron/zog-server/src/types/updates"
)

type LogController interface {
	// Insert inserts a log into the database
	Insert(ctx context.Context, log *types.Log) (err error)

	// Get gets a log from the database
	Get(ctx context.Context, id string) (log *types.Log, err error)

	// Update updates a log in the database
	Update(ctx context.Context, id string, updates *updates.Log) (err error)

	// Delete deletes a log from the database
	Delete(ctx context.Context, id string) (err error)

	// Count counts the number of logs in the database
	Count(ctx context.Context) (count int64, err error)

	// Find finds a set of logs in the database
	Find(ctx context.Context, query *find.LogQuery, options *find.LogOptions) (iter Iterator[*types.Log], err error)
}
