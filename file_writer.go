package zlog

import (
	"io"
	"os"
	"path/filepath"

	"go.innotegrity.dev/xerrors"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	// FileWriterPathError indicates there was an error with the supplied path when creating the [FileWriter].
	FileWriterPathError = 1

	// FileWriterPermissionsError indicates there was a permissions error when creating the [FileWriter].
	FileWriterPermissionsError = 2
)

// FileWriter implements a file-based writer which automatically rotates log files once they get to be
// 25MB or greater.
type FileWriter struct {
	// unexported variables
	logger *lumberjack.Logger // underlying Lumberjack logger
}

// NewFileWriter creates and initializes a new [FileWriter] object.
func NewFileWriter(file string, dirMode os.FileMode, fileMode os.FileMode, maxAge, maxCount, maxSize int) (
	*FileWriter, xerrors.Error) {

	// make sure the parent folder exists
	dir := filepath.Dir(file)
	if _, err := os.Stat(dir); err != nil {
		if !os.IsNotExist(err) {
			xerr := xerrors.Wrapf(FileWriterPathError, err, "unable to stat log folder '%s': %s", dir, err.Error()).
				WithAttrs(map[string]any{"log_file": file, "log_folder": dir})
			return nil, xerr
		}
		if err := os.MkdirAll(dir, dirMode); err != nil {
			xerr := xerrors.Wrapf(FileWriterPathError, err, "failed to create log folder '%s': %s", dir, err.Error()).
				WithAttrs(map[string]any{"log_file": file, "log_folder": dir})
			return nil, xerr
		}
	}

	// create the file if it doesn't exist ; otherwise make sure we can write to it
	h, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, fileMode)
	if err != nil {
		xerr := xerrors.Wrapf(FileWriterPathError, err, "failed to open log file '%s' for writing: %s", file, err).
			WithAttr("log_file", file)
		return nil, xerr
	}
	h.Close()

	// fix the permissions on the file
	if err := os.Chmod(file, fileMode); err != nil {
		xerr := xerrors.Wrapf(FileWriterPermissionsError, err, "failed to set permissions on log file '%s': %s",
			file, err).WithAttr("log_file", file)
		return nil, xerr
	}

	return &FileWriter{
		logger: &lumberjack.Logger{
			Compress:   false,
			Filename:   file,
			LocalTime:  false,
			MaxAge:     maxAge,
			MaxBackups: maxCount,
			MaxSize:    maxSize,
		},
	}, nil
}

// Close closes the current log file.
func (w *FileWriter) Close() error {
	if w.logger != nil {
		return w.logger.Close()
	}
	return nil
}

// Rotate forces the current log file to be rotated.
func (w *FileWriter) Rotate() error {
	if w.logger != nil {
		return w.logger.Rotate()
	}
	return nil
}

// Writer returns the underlying [io.Writer] object.
func (w *FileWriter) Writer() io.Writer {
	return w.logger
}
