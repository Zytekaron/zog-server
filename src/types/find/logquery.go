package find

import (
	"github.com/zytekaron/gotil/v2/optional"
	"github.com/zytekaron/zog-server/src/types"
	"go.mongodb.org/mongo-driver/bson"
)

type LogQuery struct {
	Services optional.Optional[[]string]      `json:"services"`
	Modules  optional.Optional[[]string]      `json:"modules"`
	Levels   optional.Optional[[]types.Level] `json:"levels"`

	Before optional.Optional[types.Time] `json:"before"`
	After  optional.Optional[types.Time] `json:"after"`
}

func (l *LogQuery) WithServices(services ...string) *LogQuery {
	l.Services = optional.Of(services)
	return l
}

func (l *LogQuery) WithModules(modules ...string) *LogQuery {
	l.Modules = optional.Of(modules)
	return l
}

func (l *LogQuery) WithLevels(levels ...types.Level) *LogQuery {
	l.Levels = optional.Of(levels)
	return l
}

func (l *LogQuery) SetBefore(before types.Time) {
	l.Before = optional.Of(before)
}

func (l *LogQuery) SetAfter(after types.Time) {
	l.After = optional.Of(after)
}

func (l *LogQuery) MongoQuery() bson.M {
	var query []bson.M
	if l.Services.IsPresent() {
		query = append(query, bson.M{"service": bson.M{"$in": l.Services.Get()}})
	}
	if l.Modules.IsPresent() {
		query = append(query, bson.M{"module": bson.M{"$in": l.Modules.Get()}})
	}
	if l.Levels.IsPresent() {
		query = append(query, bson.M{"level": bson.M{"$in": l.Levels.Get()}})
	}
	if l.Before.IsPresent() {
		query = append(query, bson.M{"created_at": bson.M{"$lt": l.Before.Get()}})
	}
	if l.After.IsPresent() {
		query = append(query, bson.M{"created_at": bson.M{"$lt": l.After.Get()}})
	}

	if len(query) == 0 {
		return bson.M{}
	}
	return bson.M{"$and": query}
}
