package zlog

import (
	"io"

	"github.com/rs/zerolog"
)

// FilteredLevelWriterConditionFunc is called to determine whether or not the given record should be logged.
type FilteredLevelWriterConditionFunc func(level zerolog.Level) bool

// FilteredLevelWriterCondition holds a single conditional function to execute.
type FilteredLevelWriterCondition struct {
	fn FilteredLevelWriterConditionFunc // function to call for the condition
}

// NewFilteredLevelWriterCondition creates a new [FilteredLevelWriterCondition] object.
func NewFilteredLevelWriterCondition(f FilteredLevelWriterConditionFunc) *FilteredLevelWriterCondition {
	return &FilteredLevelWriterCondition{
		fn: f,
	}
}

// And requires this handler condition AND the given function to be true in order to log a record.
//
// Note if either the function stored in this object or the function passed are nil, this condition
// will always return false.
func (c *FilteredLevelWriterCondition) And(f FilteredLevelWriterConditionFunc) *FilteredLevelWriterCondition {
	return &FilteredLevelWriterCondition{
		fn: func(level zerolog.Level) bool {
			if c.fn != nil && f != nil {
				return c.fn(level) && f(level)
			}
			return false
		},
	}
}

// Func returns the actual function associated with the condition that will determine whether or not to log a record.
func (c *FilteredLevelWriterCondition) Func() FilteredLevelWriterConditionFunc {
	return c.fn
}

// Or requires this handler condition OR the given function to be true in order to log a record.
//
// Note if the function stored in this object and the function passed are both nil, this condition
// will always return false.
func (c *FilteredLevelWriterCondition) Or(f FilteredLevelWriterConditionFunc) *FilteredLevelWriterCondition {
	return &FilteredLevelWriterCondition{
		fn: func(level zerolog.Level) bool {
			if c.fn != nil && f != nil {
				return c.fn(level) || f(level)
			}
			if c.fn != nil && f == nil {
				return c.fn(level)
			}
			if c.fn == nil && f != nil {
				return f(level)
			}
			return false // nil || nil = false
		},
	}
}

// FilteredLevelWriter filters messages written to the writer based on one or more conditions.
type FilteredLevelWriter struct {
	io.Writer

	// unexported variables
	cond []*FilteredLevelWriterCondition // array of conditions that must be true to write the message
}

// NewFilteredLevelWriter returns a new [FilteredLevelWriter] object.
func NewFilteredLevelWriter(w io.Writer, cond []*FilteredLevelWriterCondition) *FilteredLevelWriter {
	if w == nil {
		panic("writer should never be nil")
	}
	if cond == nil {
		panic("condition array should never be nil")
	}

	return &FilteredLevelWriter{
		Writer: w,
		cond:   cond,
	}
}

// Write writes to the underlying Writer.
func (w *FilteredLevelWriter) Write(p []byte) (int, error) {
	return w.Writer.Write(p)
}

// WriteLevel will only write the message if all of the filter conditions evaluate to true.
func (fw *FilteredLevelWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	for _, c := range fw.cond {
		if c == nil || c.fn == nil || !c.fn(level) {
			return len(p), nil
		}
	}
	return fw.Writer.Write(p)
}
