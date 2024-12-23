package logging

import (
	"fmt"
	"io"
	"os"
)

// Logger represents a simple logger
type Logger struct {
	debug    bool
	debugOut io.Writer
}

// New creates a new logger
func New(debug bool) *Logger {
	return &Logger{
		debug:    debug,
		debugOut: os.Stderr,
	}
}

// SetOutput sets the output writer for debug messages
func (l *Logger) SetOutput(w io.Writer) {
	l.debugOut = w
}

// Debug logs a debug message if debug mode is enabled
func (l *Logger) Debug(format string, args ...interface{}) {
	if l.debug {
		fmt.Fprintf(l.debugOut, format+"\n", args...)
	}
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
}
