package zlog

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

// FilteredConsoleWriter implements a [zerolog.ConsoleWriter] which sends messages with [zerolog.WarnLevel]
// and above to [os.Stderr] and all other messages to [os.Stdout].
type FilteredConsoleWriter struct {
	// unexported variables
	writers []io.Writer // writers used for writing messages
}

// NewFilteredConsoleWriter creates and initializes a new [FilteredConsoleWriter] object.
func NewFilteredConsoleWriter() *FilteredConsoleWriter {
	var stdoutWriter, stderrWriter io.Writer
	stdoutWriter = zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "03:04:05PM",
	}
	stderrWriter = zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: "03:04:05PM",
	}

	stdoutCondition := NewFilteredLevelWriterCondition(func(level zerolog.Level) bool {
		return level < zerolog.WarnLevel
	})
	stderrCondition := NewFilteredLevelWriterCondition(func(level zerolog.Level) bool {
		return level >= zerolog.WarnLevel
	})
	return &FilteredConsoleWriter{
		writers: []io.Writer{
			NewFilteredLevelWriter(stdoutWriter, []*FilteredLevelWriterCondition{stdoutCondition}),
			NewFilteredLevelWriter(stderrWriter, []*FilteredLevelWriterCondition{stderrCondition}),
		},
	}
}

// Writers returns the associated [io.Writer] objects for the writer.
func (w *FilteredConsoleWriter) Writers() []io.Writer {
	return w.writers
}
