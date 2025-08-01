// Package log provides a simplified and opinionated interface for structured logging,
// built on top of the high-performance zap logger.
package log

// Option defines a function that configures the logger.
// This follows the functional options pattern for clean and extensible configuration.
type Option func(*Config)

// WithDev configures the logger to run in development mode.
// In development mode, logs are more human-readable, with colored levels and custom time formatting.
func WithDev(development bool) Option {
	return func(cfg *Config) {
		cfg.Development = development
	}
}

// WithLevel sets the minimum logging level.
// Only logs at or above this level will be written.
// Valid levels are "debug", "info", "warn", "error".
func WithLevel(level string) Option {
	return func(cfg *Config) {
		cfg.Level = level
	}
}
