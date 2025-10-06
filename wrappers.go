package log

import (
	"log"
	"log/slog"
	"regexp"

	slogzerolog "github.com/samber/slog-zerolog/v2"
)

// WrapGoLog wraps the default go [log.Logger] object into the given [Logger].
func WrapGoLog(logger *Logger) {
	if logger == nil {
		logger = NewDiscardLogger()
	}

	log.SetOutput(newGoLogWrapper(logger))
}

// WrapGoSlog wraps the default go [slog.Logger] object into the given [Logger].
func WrapGoSlog(logger *Logger) {
	if logger == nil {
		logger = NewDiscardLogger()
	}

	l := slog.New(slogzerolog.Option{
		Level:  slog.Level(logger.GetLevel()),
		Logger: &logger.Logger,
	}.NewZerologHandler())
	slog.SetDefault(l.With("log_source", "pkg.go.dev/log/slog"))
}

// goLogWrapper wraps the standard go logger output into a [Logger] message.
type goLogWrapper struct {
	// unexported variables
	logger *Logger // internal logger to use for logging
}

// newGoLogWrapper creates and initializes a new [goLogWrapper] object.
func newGoLogWrapper(logger *Logger) *goLogWrapper {
	if logger == nil {
		logger = NewDiscardLogger()
	}

	return &goLogWrapper{
		logger: logger,
	}
}

// Write simply writes the data to the logger as a debug message.
func (l *goLogWrapper) Write(data []byte) (int, error) {
	// remove timestamp if it exists
	tsExpr := regexp.MustCompile(`^[0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} `)
	if tsExpr.Match(data) {
		data = data[20:]
	}

	l.logger.Debug().Str("log_source", "pkg.go.dev/log").Msg(string(data))
	return len(data), nil
}
