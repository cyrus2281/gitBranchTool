package internal

import (
	"fmt"
	"io"
	"os"
)

const (
	DEBUG = iota
	INFO
	WARNING
	ERROR
	FATAL
	OFF
)

type logger struct {
	level           int
	outputWriter    io.Writer
	errOutputWriter io.Writer
	prefix          string
	errPrefix       string
}

// Set the prefix for the log message
func (l *logger) SetPrefix(prefix string) {
	l.prefix = prefix
}

// Set the prefix for the error log message
func (l *logger) SetErrorPrefix(prefix string) {
	l.errPrefix = prefix
}

// Set the prefixes for the log message and the error log message
func (l *logger) SetPrefixes(prefix string, errPrefix string) {
	l.prefix = prefix
	l.errPrefix = errPrefix
}

// Set the log level
func (l *logger) SetLogLevel(level int) {
	l.level = level
}

// Get the log level
func (l *logger) GetLogLevel() int {
	return l.level
}

// Set the output writers for standard output and error output
func (l *logger) SetOutputWriters(writer io.Writer, errWriter io.Writer) {
	l.outputWriter = writer
	l.errOutputWriter = errWriter
}

// Set the output writer for standard output
func (l *logger) SetOutputWriter(writer io.Writer) {
	l.outputWriter = writer
}

// Set the output writer for error output
func (l *logger) SetErrorOutputWriter(writer io.Writer) {
	l.errOutputWriter = writer
}

// get message prefix
func (l *logger) getPrefix(level int) string {
	switch level {
	case DEBUG:
		return l.prefix + "DEBUG: "
	case INFO:
		return l.prefix + ""
	case WARNING:
		return l.prefix + "WARNING: "
	case ERROR:
		return l.errPrefix + "ERROR: "
	case FATAL:
		return l.errPrefix + "ERROR: "
	default:
		return l.prefix + ""
	}
}

func (l *logger) Log(level int, message ...any) {
	prefix := l.getPrefix(level)
	if prefix != "" {
		message = append([]any{prefix}, message...)
	}
	if l.level <= level {
		if level == ERROR {
			fmt.Fprintln(l.errOutputWriter, message...)
		} else if level == FATAL {
			fmt.Fprintln(l.errOutputWriter, message...)
			os.Exit(1)
		} else {
			fmt.Fprintln(l.outputWriter, message...)
		}
	}
}

// Log with DEBUG level
func (l *logger) Debug(message ...any) {
	l.Log(DEBUG, message...)
}

// Log with INFO level
func (l *logger) Info(message ...any) {
	l.Log(INFO, message...)
}

// Log with WARNING level
func (l *logger) Warning(message ...any) {
	l.Log(WARNING, message...)
}

// Log with ERROR level to Error output
func (l *logger) Error(message ...any) {
	l.Log(ERROR, message...)
}

// Log with FATAL level to Error output and exit with 1
func (l *logger) Fatal(message ...any) {
	l.Log(FATAL, message...)
}

// Formatted log with DEBUG level
func (l *logger) DebugF(format string, message ...any) {
	l.Log(DEBUG, fmt.Sprintf(format, message...))
}

// Formatted log with INFO level
func (l *logger) InfoF(format string, message ...any) {
	l.Log(INFO, fmt.Sprintf(format, message...))
}

// Formatted log with WARNING level
func (l *logger) WarningF(format string, message ...any) {
	l.Log(WARNING, fmt.Sprintf(format, message...))
}

// Formatted log with ERROR level to Error output
func (l *logger) ErrorF(format string, message ...any) {
	l.Log(ERROR, fmt.Sprintf(format, message...))
}

// Formatted log with FATAL level to Error output and exit with 1
func (l *logger) FatalF(format string, message ...any) {
	l.Log(FATAL, fmt.Sprintf(format, message...))
}

// Check error and log with ERROR level to Error output if not nil
func (l *logger) CheckError(err error) {
	if err != nil {
		l.Log(ERROR, err)
	}
}

// Check error and log with FATAL level to Error output and exit with 1 if not nil
func (l *logger) CheckFatal(err error) {
	if err != nil {
		l.Log(FATAL, err)
	}
}

// Check error and log with ERROR level to Error output if not nil
func (l *logger) CheckErrorF(err error, format string, message ...any) {
	if err != nil {
		l.Log(ERROR, fmt.Sprintf(format, message...), err)
	}
}

// Check error and log with FATAL level to Error output and exit with 1 if not nil
func (l *logger) CheckFatalF(err error, format string, message ...any) {
	if err != nil {
		l.Log(FATAL, fmt.Sprintf(format, message...), err)
	}
}

// Logger
var Logger logger

func init() {
	Logger = logger{
		level:           INFO,
		outputWriter:    os.Stdout,
		errOutputWriter: os.Stderr,
	}
}
