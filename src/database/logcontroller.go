package database

import (
	"context"
	"github.com/zytekaron/zog-server/src/types/find"
	"github.com/zytekaron/zog-server/src/types/updates"
)

type Controller[T any] interface {
	// Insert inserts a log into the database
	Insert(context.Context, T) error

	// Get gets a log from the database
	Get(context.Context, string) (T, error)

	// Update updates a log in the database
	Update(context.Context, string, updates.Updates[T]) error

	// Delete deletes a log from the database
	Delete(context.Context, string) error

	// Count counts the number of logs in the database
	Count(context.Context) (int64, error)

	// Find finds a set of logs in the database
	Find(context.Context, find.Query[T], find.Options[T]) (Iterator[T], error)
}
