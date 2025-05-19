package logger

import (
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
	// DisableCaller disables including the caller in log output
	DisableCaller bool
	// DisableStacktrace disables including stack traces in log output
	DisableStacktrace bool
}

// Logger represents a logger instance
type Logger struct {
	logger  *zap.Logger
	sugared *zap.SugaredLogger
}

// buildZapLogger builds a zap logger from the configuration
func buildZapLogger(cfg Config) (*zap.Logger, error) {
	// Set default output path if none provided
	if len(cfg.OutputPaths) == 0 {
		cfg.OutputPaths = []string{"stdout"}
	}

	// Parse level
	level := zap.InfoLevel
	if cfg.Level != "" {
		if err := level.Set(cfg.Level); err != nil {
			return nil, err
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
		DisableCaller:     cfg.DisableCaller,
		DisableStacktrace: cfg.DisableStacktrace,
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
		return nil, err
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

	return logger, nil
}

// Option is a function that configures the logger
type Option func(*Config)

// WithLevel sets the minimum log level
func WithLevel(level string) Option {
	return func(cfg *Config) {
		cfg.Level = level
	}
}

// WithDevelopmentMode enables or disables development mode
func WithDevelopmentMode(enabled bool) Option {
	return func(cfg *Config) {
		cfg.Development = enabled
	}
}

// WithOutputPaths sets the output paths for the logger
func WithOutputPaths(paths []string) Option {
	return func(cfg *Config) {
		cfg.OutputPaths = paths
	}
}

// WithServiceName sets the service name for the logger
func WithServiceName(name string) Option {
	return func(cfg *Config) {
		cfg.ServiceName = name
	}
}

// WithInitialFields sets initial fields for the logger
func WithInitialFields(fields map[string]interface{}) Option {
	return func(cfg *Config) {
		if cfg.InitialFields == nil {
			cfg.InitialFields = make(map[string]interface{})
		}
		for k, v := range fields {
			cfg.InitialFields[k] = v
		}
	}
}

// WithDisableCaller disables including the caller in log output
func WithDisableCaller(disable bool) Option {
	return func(cfg *Config) {
		cfg.DisableCaller = disable
	}
}

// WithDisableStacktrace disables including stack traces in log output
func WithDisableStacktrace(disable bool) Option {
	return func(cfg *Config) {
		cfg.DisableStacktrace = disable
	}
}

// NewLogger creates a new logger from the configuration
func NewLogger(cfg Config) (*Logger, error) {
	zapLogger, err := buildZapLogger(cfg)
	if err != nil {
		return nil, err
	}

	return &Logger{
		logger:  zapLogger,
		sugared: zapLogger.Sugar(),
	}, nil
}

// New creates a new logger with functional options
func New(options ...Option) (*Logger, error) {
	// Default configuration
	cfg := Config{
		Level:       "info",
		Development: false,
		OutputPaths: []string{"stdout"},
	}

	// Apply all options
	for _, option := range options {
		option(&cfg)
	}

	return NewLogger(cfg)
}

// With creates a child logger with additional fields
func (l *Logger) With(fields ...zapcore.Field) *Logger {
	// Create a new logger with the fields
	newLogger := l.logger.With(fields...)

	// Convert zapcore.Fields to interfaces for the sugared logger
	args := make([]interface{}, len(fields))
	for i, field := range fields {
		args[i] = field
	}

	return &Logger{
		logger:  newLogger,
		sugared: newLogger.Sugar().With(args...),
	}
}

// WithFields creates a child logger with additional fields as key-value pairs
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	args := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}

	return &Logger{
		logger:  l.logger,
		sugared: l.sugared.With(args...),
	}
}

// Debug logs at debug level
func (l *Logger) Debug(msg string, fields ...zapcore.Field) {
	l.logger.Debug(msg, fields...)
}

// Info logs at info level
func (l *Logger) Info(msg string, fields ...zapcore.Field) {
	l.logger.Info(msg, fields...)
}

// Warn logs at warn level
func (l *Logger) Warn(msg string, fields ...zapcore.Field) {
	l.logger.Warn(msg, fields...)
}

// Error logs at error level
func (l *Logger) Error(msg string, fields ...zapcore.Field) {
	l.logger.Error(msg, fields...)
}

// Fatal logs at fatal level and then calls os.Exit(1)
func (l *Logger) Fatal(msg string, fields ...zapcore.Field) {
	l.logger.Fatal(msg, fields...)
}

// Debugf logs at debug level with formatting (sugared logger)
func (l *Logger) Debugf(template string, args ...interface{}) {
	l.sugared.Debugf(template, args...)
}

// Debugw logs at debug level with structured key-value pairs (sugared logger)
func (l *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	l.sugared.Debugw(msg, keysAndValues...)
}

// Infof logs at info level with formatting (sugared logger)
func (l *Logger) Infof(template string, args ...interface{}) {
	l.sugared.Infof(template, args...)
}

// Infow logs at info level with structured key-value pairs (sugared logger)
func (l *Logger) Infow(msg string, keysAndValues ...interface{}) {
	l.sugared.Infow(msg, keysAndValues...)
}

// Warnf logs at warn level with formatting (sugared logger)
func (l *Logger) Warnf(template string, args ...interface{}) {
	l.sugared.Warnf(template, args...)
}

// Warnw logs at warn level with structured key-value pairs (sugared logger)
func (l *Logger) Warnw(msg string, keysAndValues ...interface{}) {
	l.sugared.Warnw(msg, keysAndValues...)
}

// Errorf logs at error level with formatting (sugared logger)
func (l *Logger) Errorf(template string, args ...interface{}) {
	l.sugared.Errorf(template, args...)
}

// Errorw logs at error level with structured key-value pairs (sugared logger)
func (l *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	l.sugared.Errorw(msg, keysAndValues...)
}

// Fatalf logs at fatal level with formatting and then calls os.Exit(1)
func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.sugared.Fatalf(template, args...)
}

// Fatalw logs at fatal level with structured key-value pairs (sugared logger) and then calls os.Exit(1)
func (l *Logger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.sugared.Fatalw(msg, keysAndValues...)
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() error {
	return l.logger.Sync()
}

// NewNopLogger returns a no-op logger for testing where logs are undesired
func NewNopLogger() *Logger {
	// Create a no-op zap logger
	noopCore := zapcore.NewNopCore()
	noopZap := zap.New(noopCore)
	
	return &Logger{
		logger:  noopZap,
		sugared: noopZap.Sugar(),
	}
}

// NoOp returns a no-op logger for testing where logs are undesired
// This is an alias for NewNopLogger for backward compatibility
func NoOp() *Logger {
	return NewNopLogger()
}
