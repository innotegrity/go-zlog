package zlog

import (
	"encoding/json"

	"github.com/rs/zerolog"
)

// [zerolog.Level] mappings
const (
	TraceLevel = Level(zerolog.TraceLevel)
	DebugLevel = Level(zerolog.DebugLevel)
	InfoLevel  = Level(zerolog.InfoLevel)
	WarnLevel  = Level(zerolog.WarnLevel)
	ErrorLevel = Level(zerolog.ErrorLevel)
	FatalLevel = Level(zerolog.FatalLevel)
	PanicLevel = Level(zerolog.PanicLevel)
	NoLevel    = Level(zerolog.NoLevel)
	Disabled   = Level(zerolog.Disabled)
)

// Level represents a log level.
type Level zerolog.Level

// ParseLevel parses the given string into a proper [Level] object.
func ParseLevel(l string) (Level, error) {
	level, err := zerolog.ParseLevel(l)
	if err != nil {
		return NoLevel, err
	}
	return Level(level), nil
}

// MarshalJSON marshals the [Level] object to JSON.
func (l Level) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.String())
}

// MarshalText marshals the [Level] object to plain text.
func (l Level) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

// String returns the [Level] object as a string.
func (l Level) String() string {
	return zerolog.Level(l).String()
}

// UnmarshalJSON parses the JSON data into a [Level] object.
func (l *Level) UnmarshalJSON(data []byte) error {
	var lvl string
	if err := json.Unmarshal(data, &lvl); err != nil {
		return err
	}
	parsedLvl, err := zerolog.ParseLevel(lvl)
	if err != nil {
		return err
	}
	*l = Level(parsedLvl)
	return nil
}

// UnmarshalText parses the text into a [Level] object.
func (l *Level) UnmarshalText(data []byte) error {
	parsedLvl, err := zerolog.ParseLevel(string(data))
	if err != nil {
		return err
	}
	*l = Level(parsedLvl)
	return nil
}

// Z returns the equivalent [zerolog.Level].
func (l Level) Z() zerolog.Level {
	return zerolog.Level(l)
}
