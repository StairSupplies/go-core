package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Level type
type Level = zapcore.Level

// Log levels
const (
	DebugLevel = zapcore.DebugLevel
	InfoLevel  = zapcore.InfoLevel
	WarnLevel  = zapcore.WarnLevel
	ErrorLevel = zapcore.ErrorLevel
	FatalLevel = zapcore.FatalLevel
)

// Config contains configuration for the logger
type Config struct {
	// Level is the minimum log level to output
	Level string
	// Development sets development mode (more human-friendly output)
	Development bool
	// OutputPaths defines where logs are written to
	OutputPaths []string
	// ServiceName is the name of the service for all logs
	ServiceName string
	// InitialFields are fields added to all log entries
	InitialFields map[string]interface{}
}

var (
	// global logger instance
	globalLogger *zap.Logger
	// global sugared logger (easier to use)
	globalSugared *zap.SugaredLogger
)

type ctxLoggerKey struct{}

// Init initializes the global logger based on the provided config
func Init(cfg Config) error {
	// Set default output path if none provided
	if len(cfg.OutputPaths) == 0 {
		cfg.OutputPaths = []string{"stdout"}
	}

	// Parse level
	level := zap.InfoLevel
	if cfg.Level != "" {
		if err := level.Set(cfg.Level); err != nil {
			return err
		}
	}

	// Create zap config
	zapConfig := zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		Development:       cfg.Development,
		Encoding:          "json",
		EncoderConfig:     zap.NewProductionEncoderConfig(),
		OutputPaths:       cfg.OutputPaths,
		ErrorOutputPaths:  []string{"stderr"},
		DisableCaller:     false,
		DisableStacktrace: false,
	}

	// Use more human-friendly settings for development
	if cfg.Development {
		zapConfig.Encoding = "console"
		zapConfig.EncoderConfig = zap.NewDevelopmentEncoderConfig()
	} else {
		// Production settings
		zapConfig.EncoderConfig.TimeKey = "timestamp"
		zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	// Build the logger
	logger, err := zapConfig.Build(zap.AddCallerSkip(1))
	if err != nil {
		return err
	}

	// Add default fields
	fields := []zap.Field{}
	if cfg.ServiceName != "" {
		fields = append(fields, zap.String("service", cfg.ServiceName))
	}

	for k, v := range cfg.InitialFields {
		fields = append(fields, zap.Any(k, v))
	}

	if len(fields) > 0 {
		logger = logger.With(fields...)
	}

	// Set the global logger instances
	globalLogger = logger
	globalSugared = logger.Sugar()

	return nil
}

// GetLogger returns the global logger
func GetLogger() *zap.Logger {
	if globalLogger == nil {
		// Create a default logger if not initialized
		globalLogger, _ = zap.NewProduction()
		globalSugared = globalLogger.Sugar()
	}
	return globalLogger
}

// GetSugared returns the global sugared logger
func GetSugared() *zap.SugaredLogger {
	if globalSugared == nil {
		// Create a default logger if not initialized
		globalLogger, _ = zap.NewProduction()
		globalSugared = globalLogger.Sugar()
	}
	return globalSugared
}

// WithContext returns a logger from the context or the default logger
func WithContext(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(ctxLoggerKey{}).(*zap.Logger); ok {
		return l
	}
	return GetLogger()
}

// ContextWithLogger adds a logger to the context
func ContextWithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxLoggerKey{}, logger)
}

// With creates a child logger with additional fields
func With(fields ...zapcore.Field) *zap.Logger {
	return GetLogger().With(fields...)
}

// WithFields creates a child sugared logger with additional fields
func WithFields(fields map[string]interface{}) *zap.SugaredLogger {
	args := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}
	return GetSugared().With(args...)
}

// Debug logs at debug level
func Debug(msg string, fields ...zapcore.Field) {
	GetLogger().Debug(msg, fields...)
}

// Info logs at info level
func Info(msg string, fields ...zapcore.Field) {
	GetLogger().Info(msg, fields...)
}

// Warn logs at warn level
func Warn(msg string, fields ...zapcore.Field) {
	GetLogger().Warn(msg, fields...)
}

// Error logs at error level
func Error(msg string, fields ...zapcore.Field) {
	GetLogger().Error(msg, fields...)
}

// Fatal logs at fatal level and then calls os.Exit(1)
func Fatal(msg string, fields ...zapcore.Field) {
	GetLogger().Fatal(msg, fields...)
}

// Debugf logs at debug level with formatting (sugared logger)
func Debugf(template string, args ...interface{}) {
	GetSugared().Debugf(template, args...)
}

// Infof logs at info level with formatting (sugared logger)
func Infof(template string, args ...interface{}) {
	GetSugared().Infof(template, args...)
}

// Warnf logs at warn level with formatting (sugared logger)
func Warnf(template string, args ...interface{}) {
	GetSugared().Warnf(template, args...)
}

// Errorf logs at error level with formatting (sugared logger)
func Errorf(template string, args ...interface{}) {
	GetSugared().Errorf(template, args...)
}

// Fatalf logs at fatal level with formatting and then calls os.Exit(1)
func Fatalf(template string, args ...interface{}) {
	GetSugared().Fatalf(template, args...)
}

// Sync flushes any buffered log entries
func Sync() error {
	return GetLogger().Sync()
}
