package log

import (
	"context"
	"io"

	"github.com/rs/zerolog"
)

var (
	// CallerSkipFrameCount is the number of stack frames to skip to find the caller.
	CallerSkipFrameCount = zerolog.CallerSkipFrameCount
)

// Logger wraps a [zerolog.Logger] object, saving its associated [io.Writer] objects for later retrieval.
type Logger struct {
	zerolog.Logger

	// unexported variables
	writers []io.Writer // writers that the zerolog logger is using
}

// FromContext retrieves the logger from the given context.
//
// If no logger is found in the context, [NewDiscardLogger] is called to return a logger that discards all messages.
func FromContext(ctx context.Context) *Logger {
	if v := ctx.Value(loggerCtxKey{}); v != nil {
		if logger, ok := v.(*Logger); ok {
			return logger
		}
	}
	return NewDiscardLogger()
}

// NewDiscardLogger creates and initializes a [Logger] object that simply discards all messages.
func NewDiscardLogger() *Logger {
	return &Logger{
		Logger:  zerolog.New(io.Discard).Level(zerolog.Disabled),
		writers: []io.Writer{io.Discard},
	}
}

// NewLogger creates and initializes a new [Logger] object.
//
// If includeCaller is true or the level is set at [zerolog.DebugLevel] or lower, caller information is automatically
// included in log messages.
func NewLogger(logLevel Level, includeCaller bool, writers ...io.Writer) *Logger {
	if len(writers) == 0 {
		writers = []io.Writer{io.Discard}
	}

	multiWriter := zerolog.MultiLevelWriter(writers...)
	l := &Logger{
		Logger:  zerolog.New(multiWriter).With().Timestamp().Logger().Level(zerolog.Level(logLevel)),
		writers: writers,
	}

	if includeCaller || logLevel <= Level(zerolog.DebugLevel) {
		l.Logger = l.Logger.With().Caller().Logger()
	}
	return l
}

// Close closes any writers that implement [io.Closer].
func (l *Logger) Close() {
	for _, writer := range l.writers {
		if w, ok := writer.(io.Closer); ok {
			w.Close()
		}
	}
}

// IsDebugEnabled returns whether or not debug logging is enabled.
func (l *Logger) IsDebugEnabled() bool {
	return l.GetLevel() <= zerolog.DebugLevel
}

// ReplaceLevel replaces the minimum logging level for the underlying [zerolog.Logger] and returns the previous level.
func (l *Logger) ReplaceLevel(logLevel Level) Level {
	oldLevel := l.Logger.GetLevel()
	l.Logger = l.Logger.Level(zerolog.Level(logLevel))
	return Level(oldLevel)
}

// WithContext attaches this logger to the given context and returns the new context.
func (l *Logger) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, loggerCtxKey{}, l)
}

// WithZ creates a new child logger using the given [zerolog.Logger] and the current object writers and returns
// the new logger.
func (l *Logger) WithZ(logger zerolog.Logger) *Logger {
	return &Logger{
		Logger:  logger,
		writers: l.writers,
	}
}

// WithZContext creates a new child logger with the given [zerolog.Context] and the current object writers and returns
// the new logger.
func (l *Logger) WithZContext(ctx zerolog.Context) *Logger {
	return &Logger{
		Logger:  ctx.Logger(),
		writers: l.writers,
	}
}

// Writers returns the underlying [io.Writer] objects used by the logger.
func (l *Logger) Writers() []io.Writer {
	return l.writers
}

// Z returns the underlying [zerolog.Logger] object.
func (l *Logger) Z() *zerolog.Logger {
	return &l.Logger
}

// loggerCtxKey is an empty struct used to add the logger to context.
type loggerCtxKey struct{}
