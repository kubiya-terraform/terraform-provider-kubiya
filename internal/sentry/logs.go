package sentry

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
)

// LogLevel represents the severity level of a log message
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)

// Logger is the interface for Sentry-integrated logging
type Logger interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Fatal(msg string, fields ...interface{})
	WithContext(ctx context.Context) Logger
	WithFields(fields map[string]interface{}) Logger
}

// SentryLogger implements Logger with Sentry integration using slog
type SentryLogger struct {
	ctx        context.Context
	slogLogger *slog.Logger
}

// NewLogger creates a new Sentry-integrated logger using slog
func NewLogger() Logger {
	handler := NewSentryHandler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

	return &SentryLogger{
		ctx:        context.Background(),
		slogLogger: slog.New(handler),
	}
}

// SentryHandler wraps slog.Handler to send logs to Sentry
type SentryHandler struct {
	inner slog.Handler
}

// NewSentryHandler creates a new Sentry handler for slog
func NewSentryHandler(inner slog.Handler) *SentryHandler {
	return &SentryHandler{inner: inner}
}

// Enabled reports whether the handler handles records at the given level
func (h *SentryHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.inner.Enabled(ctx, level)
}

// Handle handles the Record, sending appropriate logs to Sentry
func (h *SentryHandler) Handle(ctx context.Context, record slog.Record) error {
	// First, let the inner handler process the record
	if err := h.inner.Handle(ctx, record); err != nil {
		return err
	}

	// Send to Sentry based on level
	switch record.Level {
	case slog.LevelError:
		h.sendToSentry(ctx, record, sentry.LevelError)
	case slog.LevelWarn:
		h.sendToSentry(ctx, record, sentry.LevelWarning)
	case slog.LevelInfo:
		h.sendToSentry(ctx, record, sentry.LevelInfo)
	case slog.LevelDebug:
		h.sendToSentry(ctx, record, sentry.LevelDebug)
	}

	return nil
}

// WithAttrs returns a new Handler with additional attributes
func (h *SentryHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &SentryHandler{inner: h.inner.WithAttrs(attrs)}
}

// WithGroup returns a new Handler with a group name
func (h *SentryHandler) WithGroup(name string) slog.Handler {
	return &SentryHandler{inner: h.inner.WithGroup(name)}
}

// sendToSentry sends a log record to Sentry
func (h *SentryHandler) sendToSentry(ctx context.Context, record slog.Record, level sentry.Level) {
	// Extract attributes from the record
	attrs := make(map[string]interface{})
	record.Attrs(func(a slog.Attr) bool {
		attrs[a.Key] = a.Value.Any()
		return true
	})

	// Add source information
	if record.PC != 0 {
		frames := runtime.CallersFrames([]uintptr{record.PC})
		if frame, _ := frames.Next(); frame.PC != 0 {
			attrs["source.file"] = frame.File
			attrs["source.line"] = frame.Line
			attrs["source.function"] = frame.Function
		}
	}

	// Create breadcrumb for all log levels
	AddBreadcrumb("log", record.Message, level, attrs)

	// For error and fatal levels, also capture as an event
	if level >= sentry.LevelError {
		sentry.WithScope(func(scope *sentry.Scope) {
			scope.SetLevel(level)
			scope.SetContext("log", attrs)

			// Add trace ID if available
			if traceID := GetTraceID(ctx); traceID != "" {
				scope.SetTag("trace_id", traceID)
			}

			// Check if there's an error in the attributes
			if err, ok := attrs["error"]; ok {
				if e, ok := err.(error); ok {
					sentry.CaptureException(e)
				} else {
					sentry.CaptureMessage(record.Message)
				}
			} else {
				sentry.CaptureMessage(record.Message)
			}
		})
	}
}

// Debug logs a debug message
func (l *SentryLogger) Debug(msg string, fields ...interface{}) {
	l.logWithLevel(LogLevelDebug, msg, fields...)
}

// Info logs an info message
func (l *SentryLogger) Info(msg string, fields ...interface{}) {
	l.logWithLevel(LogLevelInfo, msg, fields...)
}

// Warn logs a warning message
func (l *SentryLogger) Warn(msg string, fields ...interface{}) {
	l.logWithLevel(LogLevelWarn, msg, fields...)
}

// Error logs an error message
func (l *SentryLogger) Error(msg string, fields ...interface{}) {
	l.logWithLevel(LogLevelError, msg, fields...)
}

// Fatal logs a fatal message and exits
func (l *SentryLogger) Fatal(msg string, fields ...interface{}) {
	l.logWithLevel(LogLevelFatal, msg, fields...)
	Flush() // Ensure all events are sent before exit
	os.Exit(1)
}

// WithContext returns a new logger with context
func (l *SentryLogger) WithContext(ctx context.Context) Logger {
	return &SentryLogger{
		ctx:        ctx,
		slogLogger: l.slogLogger.With("trace_id", GetTraceID(ctx)),
	}
}

// WithFields returns a new logger with additional fields
func (l *SentryLogger) WithFields(fields map[string]interface{}) Logger {
	// Convert fields to slog attributes
	attrs := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		attrs = append(attrs, k, v)
	}

	return &SentryLogger{
		ctx:        l.ctx,
		slogLogger: l.slogLogger.With(attrs...),
	}
}

// logWithLevel logs a message with the specified level using slog
func (l *SentryLogger) logWithLevel(level LogLevel, msg string, fields ...interface{}) {
	// Convert fields to slog args format (alternating key-value pairs)
	args := make([]any, 0, len(fields))
	for i := 0; i < len(fields)-1; i += 2 {
		if key, ok := fields[i].(string); ok {
			args = append(args, key, fields[i+1])
		}
	}

	// Log using slog with appropriate level
	switch level {
	case LogLevelDebug:
		l.slogLogger.DebugContext(l.ctx, msg, args...)
	case LogLevelInfo:
		l.slogLogger.InfoContext(l.ctx, msg, args...)
	case LogLevelWarn:
		l.slogLogger.WarnContext(l.ctx, msg, args...)
	case LogLevelError:
		l.slogLogger.ErrorContext(l.ctx, msg, args...)
	case LogLevelFatal:
		l.slogLogger.ErrorContext(l.ctx, fmt.Sprintf("FATAL: %s", msg), args...)
	}
}

// StandardLoggerHook redirects standard library log output to Sentry
type StandardLoggerHook struct {
	logger Logger
}

// NewStandardLoggerHook creates a hook for standard library logger
func NewStandardLoggerHook(logger Logger) *StandardLoggerHook {
	return &StandardLoggerHook{logger: logger}
}

// Write implements io.Writer interface to capture standard log output
func (h *StandardLoggerHook) Write(p []byte) (n int, err error) {
	msg := string(p)
	msg = strings.TrimSpace(msg)

	// Detect log level from message prefix
	switch {
	case strings.HasPrefix(msg, "ERROR:"), strings.HasPrefix(msg, "[ERROR]"):
		h.logger.Error(msg)
	case strings.HasPrefix(msg, "WARN:"), strings.HasPrefix(msg, "[WARN]"):
		h.logger.Warn(msg)
	case strings.HasPrefix(msg, "INFO:"), strings.HasPrefix(msg, "[INFO]"):
		h.logger.Info(msg)
	case strings.HasPrefix(msg, "DEBUG:"), strings.HasPrefix(msg, "[DEBUG]"):
		h.logger.Debug(msg)
	default:
		h.logger.Info(msg)
	}

	return len(p), nil
}

// ConfigureStandardLogger redirects standard library logger to Sentry
func ConfigureStandardLogger(logger Logger) {
	hook := NewStandardLoggerHook(logger)
	log.SetOutput(hook)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// LoggerFromContext retrieves a logger from context or returns default
func LoggerFromContext(ctx context.Context) Logger {
	if logger, ok := ctx.Value(contextKeyLogger).(Logger); ok {
		return logger
	}
	return defaultLogger
}

// ContextWithLogger adds a logger to context
func ContextWithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, contextKeyLogger, logger)
}

// Package-level logger instance
var defaultLogger Logger

// Package-level context key for logger
const contextKeyLogger contextKey = "sentry_logger"

// InitializeLogging sets up the logging system
func InitializeLogging() {
	defaultLogger = NewLogger()
	ConfigureStandardLogger(defaultLogger)
}

// GetLogger returns the default logger instance
func GetLogger() Logger {
	if defaultLogger == nil {
		InitializeLogging()
	}
	return defaultLogger
}

// LogResourceOperation logs a resource operation with structured data
func LogResourceOperation(ctx context.Context, operation, resourceType, resourceID string, fields map[string]interface{}) {
	logger := LoggerFromContext(ctx).WithContext(ctx)

	// Merge provided fields with standard fields
	allFields := map[string]interface{}{
		"operation":     operation,
		"resource_type": resourceType,
		"resource_id":   resourceID,
		"timestamp":     time.Now().Unix(),
	}

	for k, v := range fields {
		allFields[k] = v
	}

	logger = logger.WithFields(allFields)
	logger.Info(fmt.Sprintf("Resource operation: %s %s", operation, resourceType))
}

// LogAPICall logs an API call with structured data
func LogAPICall(ctx context.Context, method, url string, statusCode int, duration time.Duration) {
	logger := LoggerFromContext(ctx).WithContext(ctx)

	fields := map[string]interface{}{
		"http.method":      method,
		"http.url":         url,
		"http.status_code": statusCode,
		"duration_ms":      duration.Milliseconds(),
	}

	logger = logger.WithFields(fields)

	if statusCode >= 400 {
		logger.Error(fmt.Sprintf("API call failed: %s %s returned %d", method, url, statusCode))
	} else {
		logger.Info(fmt.Sprintf("API call: %s %s", method, url))
	}
}
