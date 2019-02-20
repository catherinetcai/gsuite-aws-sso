package logging

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// The global logger
	logger *zap.Logger
)

// Configure configures the zap logger for Stackdriver
func Configure(env string) (err error) {
	if env == "" {
		fmt.Println("using development logger")
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger, err = config.Build()
		return
	}

	fmt.Println("using production logger")
	logger, err = zap.NewProduction()
	return
}

// Logger returns the global instance of a logger
func Logger() *zap.Logger {
	if logger == nil {
		Configure("")
	}
	return logger
}

// SetLogger allows the overriding of the global logger, really
// only recommended for testing.
func SetLogger(l *zap.Logger) {
	logger = l
}
