package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type Time time.Time

func (j Time) MarshalJSON() ([]byte, error) {
	ms := time.Time(j).UnixMilli()
	i := strconv.FormatInt(ms, 10)
	return json.Marshal(i)
}

func (j *Time) UnmarshalJSON(data []byte) error {
	var ms int64
	err := json.Unmarshal(data, &ms)
	if err != nil {
		return err
	}

	*j = Time(time.Unix(0, ms*int64(time.Millisecond)))
	return nil
}

func (j Time) MarshalBSONValue() (bsontype.Type, []byte, error) {
	ms := time.Time(j).UnixMilli()
	return bson.MarshalValue(ms)
}

func (j *Time) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.Int64 {
		return fmt.Errorf("invalid bson value type '%s'", t.String())
	}
	ms, _, ok := bsoncore.ReadInt64(data)
	if !ok {
		return errors.New("invalid bson string value")
	}

	*j = Time(time.Unix(0, ms*int64(time.Millisecond)))
	return nil
}

type Duration time.Duration

func (j Duration) MarshalJSON() ([]byte, error) {
	ms := time.Duration(j).Milliseconds()
	i := strconv.FormatInt(ms, 10)
	return json.Marshal(i)
}

func (j *Duration) UnmarshalJSON(data []byte) error {
	var ms int64
	err := json.Unmarshal(data, &ms)
	if err != nil {
		return err
	}

	*j = Duration(ms * int64(time.Millisecond))
	return nil
}

func (j Duration) MarshalBSONValue() (bsontype.Type, []byte, error) {
	ms := time.Duration(j).Milliseconds()
	return bson.MarshalValue(ms)
}

func (j *Duration) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.Int64 {
		return fmt.Errorf("invalid bson value type '%s'", t.String())
	}
	ms, _, ok := bsoncore.ReadInt64(data)
	if !ok {
		return errors.New("invalid bson string value")
	}

	*j = Duration(time.Duration(ms) * time.Millisecond)
	return nil
}
