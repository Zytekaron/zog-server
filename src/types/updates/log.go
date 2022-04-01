package updates

import (
	"github.com/zytekaron/gotil/v2/optional"
	"github.com/zytekaron/zog-server/src/types"
	"gopkg.in/mgo.v2/bson"
)

type Log struct {
	Level     optional.Optional[types.Level]
	Service   optional.Optional[string]
	Module    optional.Optional[string]
	Message   optional.Optional[string]
	CreatedAt optional.Optional[types.Time]
}

func (l *Log) WithLevel(level types.Level) *Log {
	l.Level = optional.Of(level)
	return l
}

func (l *Log) WithService(service string) *Log {
	l.Service = optional.Of(service)
	return l
}

func (l *Log) WithModule(module string) *Log {
	l.Module = optional.Of(module)
	return l
}

func (l *Log) WithMessage(message string) *Log {
	l.Message = optional.Of(message)
	return l
}

func (l *Log) WithCreatedAt(createdAt types.Time) *Log {
	l.CreatedAt = optional.Of(createdAt)
	return l
}

func (l *Log) Apply(o *types.Log) {
	o.Level = l.Level.OrElse(o.Level)
	o.Service = l.Service.OrElse(o.Service)
	o.Module = l.Module.OrElse(o.Module)
	o.Message = l.Message.OrElse(o.Message)
	o.CreatedAt = l.CreatedAt.OrElse(o.CreatedAt)
}

func (l *Log) MongoQuery() bson.M {
	set := bson.M{}
	if l.Level.IsPresent() {
		set["level"] = l.Level.Get()
	}
	if l.Service.IsPresent() {
		set["service"] = l.Service.Get()
	}
	if l.Module.IsPresent() {
		set["module"] = l.Module.Get()
	}
	if l.Message.IsPresent() {
		set["message"] = l.Message.Get()
	}
	if l.CreatedAt.IsPresent() {
		set["created_at"] = l.CreatedAt.Get()
	}

	return bson.M{
		"$set": set,
	}
}
