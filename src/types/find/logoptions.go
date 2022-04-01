package find

import (
	"github.com/zytekaron/gotil/v2/optional"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LogOptions struct {
	skip  optional.Optional[int64]
	limit optional.Optional[int64]
	sort  optional.Optional[logSort]
}

type logSort struct {
	Field string
	Order int
}

func NewLogOptions() *LogOptions {
	return &LogOptions{}
}

func (l *LogOptions) Skip(skip int64) *LogOptions {
	l.skip = optional.Of(skip)
	return l
}

func (l *LogOptions) Limit(limit int64) *LogOptions {
	l.limit = optional.Of(limit)
	return l
}

func (l *LogOptions) SortCreatedAt(order int) *LogOptions {
	l.sort = optional.Of(logSort{
		Field: "created_at",
		Order: order,
	})
	return l
}

func (l *LogOptions) MongoOptions() *options.FindOptions {
	opts := options.Find()
	if l.skip.IsPresent() {
		opts.SetSkip(l.skip.Get())
	}
	if l.limit.IsPresent() {
		opts.SetLimit(l.limit.Get())
	}
	if l.sort.IsPresent() {
		s := l.sort.Get()
		opts.SetSort(bson.M{s.Field: s.Order})
	}
	return opts
}
