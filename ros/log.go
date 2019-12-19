package ros

import (
	"fmt"
	"log"
	"os"
)

// LogLevel defines the logging level of default rosgo logger.
type LogLevel int

const (
	//LogLevelDebug is for logging debugging entries that are very verbose.
	LogLevelDebug LogLevel = iota
	// LogLevelInfo is for logging information entries that gives an idea of what's going on in the application.
	LogLevelInfo
	// LogLevelWarn is for logging non-critical entries that needs to be taken a look at.
	LogLevelWarn
	// LogLevelError is for logging error entries that represent something failing in the application.
	LogLevelError
	// LogLevelFatal is for logging unrecoverable failures and exiting (non-zero) the application.
	LogLevelFatal
)

// Logger defines an interface that a ros logger should implement
type Logger interface {
	Severity() LogLevel
	SetSeverity(severity LogLevel)
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
}

// DefaultLogger is the default logger that rosgo uses to log various logging
// entries. DefaultLogger implements the logger interface that rosgo defines.
type DefaultLogger struct {
	severity LogLevel
}

// NewDefaultLogger creates and returns a new rosgo logger instance
func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{LogLevelInfo}
}

// Severity returns the logging level of DefaultLogger
func (logger *DefaultLogger) Severity() LogLevel {
	return logger.severity
}

// SetSeverity sets the logging level of DefaultLogger
func (logger *DefaultLogger) SetSeverity(severity LogLevel) {
	logger.severity = severity
}

// Debug logs entry/entries v when the logger level is set to DebugLevel or less
func (logger *DefaultLogger) Debug(v ...interface{}) {
	if int(logger.severity) <= int(LogLevelDebug) {
		msg := fmt.Sprintf("[DEBUG] %s", fmt.Sprint(v...))
		log.Println(msg)
	}
}

// Debugf formats and logs entry/entries v when the logger level is set to DebugLevel or less
func (logger *DefaultLogger) Debugf(format string, v ...interface{}) {
	if int(logger.severity) <= int(LogLevelDebug) {
		log.Printf("[DEBUG] "+format, v...)
	}
}

// Info logs entry/entries v when the logger level is set to InfoLevel or less
func (logger *DefaultLogger) Info(v ...interface{}) {
	if int(logger.severity) <= int(LogLevelInfo) {
		msg := fmt.Sprintf("[INFO] %s", fmt.Sprint(v...))
		log.Println(msg)
	}
}

// Infof formats and logs entry/entries v when the logger level is set to InfoLevel or less
func (logger *DefaultLogger) Infof(format string, v ...interface{}) {
	if int(logger.severity) <= int(LogLevelInfo) {
		log.Printf("[INFO] "+format, v...)
	}
}

// Warn logs entry/entries v when the logger level is set to WarnLevel or less
func (logger *DefaultLogger) Warn(v ...interface{}) {
	if int(logger.severity) <= int(LogLevelWarn) {
		msg := fmt.Sprintf("[WARN] %s", fmt.Sprint(v...))
		log.Println(msg)
	}
}

// Warnf formats and logs entry/entries v when the logger level is set to WarnLevel or less
func (logger *DefaultLogger) Warnf(format string, v ...interface{}) {
	if int(logger.severity) <= int(LogLevelWarn) {
		log.Printf("[WARN] "+format, v...)
	}
}

// Error logs entry/entries v when the logger level is set to ErrorLevel or less
func (logger *DefaultLogger) Error(v ...interface{}) {
	if int(logger.severity) <= int(LogLevelError) {
		msg := fmt.Sprintf("[ERROR] %s", fmt.Sprint(v...))
		log.Println(msg)
	}
}

// Errorf formats and logs entry/entries v when the logger level is set to ErrorLevel or less
func (logger *DefaultLogger) Errorf(format string, v ...interface{}) {
	if int(logger.severity) <= int(LogLevelError) {
		log.Printf("[ERROR] "+format, v...)
	}
}

// Fatal logs entry/entries v and exits the program when the logger level is set to FatalLevel or less
func (logger *DefaultLogger) Fatal(v ...interface{}) {
	if int(logger.severity) <= int(LogLevelFatal) {
		msg := fmt.Sprintf("[FATAL] %s", fmt.Sprint(v...))
		log.Println(msg)
		os.Exit(1)
	}
}

// Fatalf formats and logs entry/entries v and exits the program when the logger level is set to FatalLevel or less
func (logger *DefaultLogger) Fatalf(format string, v ...interface{}) {
	if int(logger.severity) <= int(LogLevelFatal) {
		log.Printf("[FATAL] "+format, v...)
		os.Exit(1)
	}
}
