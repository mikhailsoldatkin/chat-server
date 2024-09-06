package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger
var sugarLogger *zap.SugaredLogger

// Init initializes the global logger with the provided zap core and options.
// It also sets up the sugared logger for use with formatted logging.
func Init(core zapcore.Core, options ...zap.Option) {
	globalLogger = zap.New(core, options...)
	sugarLogger = globalLogger.Sugar()
}

func Logger() *zap.Logger {
	return globalLogger
}

// Debug logs a debug message with structured fields.
func Debug(msg string, fields ...zap.Field) {
	globalLogger.Debug(msg, fields...)
}

// Debugf logs a debug message using a formatted string.
func Debugf(msg string, values ...any) {
	sugarLogger.Debugf(msg, values...)
}

// Info logs an informational message with structured fields.
func Info(msg string, fields ...zap.Field) {
	globalLogger.Info(msg, fields...)
}

// Infof logs an informational message using a formatted string.
func Infof(msg string, values ...any) {
	sugarLogger.Infof(msg, values...)
}

// Warn logs a warning message with structured fields.
func Warn(msg string, fields ...zap.Field) {
	globalLogger.Warn(msg, fields...)
}

// Warnf logs a warning message using a formatted string.
func Warnf(msg string, values ...any) {
	sugarLogger.Warnf(msg, values...)
}

// Error logs an error message with structured fields.
func Error(msg string, fields ...zap.Field) {
	globalLogger.Error(msg, fields...)
}

// Errorf logs an error message using a formatted string.
func Errorf(msg string, values ...any) {
	sugarLogger.Errorf(msg, values...)
}

// Fatal logs a fatal message with structured fields, then exits the application.
func Fatal(msg string, fields ...zap.Field) {
	globalLogger.Fatal(msg, fields...)
}

// Fatalf logs a fatal message using a formatted string, then exits the application.
func Fatalf(msg string, values ...any) {
	sugarLogger.Fatalf(msg, values...)
}

// WithOptions clones the global logger and applies the supplied options.
// It returns a new *zap.Logger instance with the specified options.
func WithOptions(opts ...zap.Option) *zap.Logger {
	return globalLogger.WithOptions(opts...)
}
