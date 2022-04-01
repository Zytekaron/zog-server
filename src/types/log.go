package types

import (
	"fmt"
	"time"
)

// LogIDChars calculation: https://zelark.github.io/nano-id-cc/
const LogIDChars = "0123456789abcdef"

// LogIDLength calculation: https://zelark.github.io/nano-id-cc/
const LogIDLength = 24

// LogAllowedTimeDrift is the maximum amount of time that the
// Log.CreatedAt value can be in the future relative to the
// current system clock. A low value is recommended to prevent
// future dates from being accepted for insertion into the db
const LogAllowedTimeDrift = 5 * time.Minute

const LogMaxServiceLength = 64
const LogMaxModuleLength = 64
const LogMaxMessageLength = 4 * 1024

// efficient lookup for Log.Errors
var idCharMap = map[rune]struct{}{}

func init() {
	for _, c := range []rune(LogIDChars) {
		idCharMap[c] = struct{}{}
	}
}

type Log struct {
	ID        string `json:"id,omitempty" bson:"_id"`
	Level     Level  `json:"level,omitempty" bson:"level"`
	Service   string `json:"service,omitempty" bson:"service"`
	Module    string `json:"module,omitempty" bson:"module"`
	Message   string `json:"message,omitempty" bson:"message"`
	CreatedAt Time   `json:"created_at,omitempty" bson:"created_at"`
}

type Level string

const (
	LevelFatal Level = "FATAL"
	LevelError Level = "ERROR"
	LevelWarn  Level = "WARN"
	LevelInfo  Level = "INFO"
	LevelDebug Level = "DEBUG"
	LevelTrace Level = "TRACE"
)

var LogLevels = map[Level]int{
	LevelFatal: 6,
	LevelError: 5,
	LevelWarn:  4,
	LevelInfo:  3,
	LevelDebug: 2,
	LevelTrace: 1,
}

func (l *Log) Errors() []string {
	var errors []string

	if len(l.ID) != LogIDLength {
		err := fmt.Sprintf("id must be exactly %d bytes", LogIDLength)
		errors = append(errors, err)
	}
	for _, c := range []rune(l.ID) {
		if i, ok := idCharMap[c]; !ok {
			err := fmt.Sprintf("invalid byte in id at position %d", i)
			errors = append(errors, err)
			break
		}
	}

	if _, ok := LogLevels[l.Level]; !ok {
		errors = append(errors, "level is not valid")
	}

	if len(l.Service) == 0 || len(l.Service) > LogMaxServiceLength {
		err := fmt.Sprintf("service must be 1 to %d bytes", LogMaxServiceLength)
		errors = append(errors, err)
	}

	if len(l.Module) == 0 || len(l.Module) > LogMaxModuleLength {
		err := fmt.Sprintf("module must be 1 to %d bytes", LogMaxModuleLength)
		errors = append(errors, err)
	}

	if len(l.Message) == 0 || len(l.Message) > LogMaxMessageLength {
		err := fmt.Sprintf("message must be 1 to %d bytes", LogMaxMessageLength)
		errors = append(errors, err)
	}

	t := time.Time(l.CreatedAt)
	if t.IsZero() {
		errors = append(errors, "created_at is zero")
	}
	if t.Sub(time.Now()).Nanoseconds() > LogAllowedTimeDrift.Nanoseconds() {
		errors = append(errors, "created_at is in the future")
	}

	return errors
}
