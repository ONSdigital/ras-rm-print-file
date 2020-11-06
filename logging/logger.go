package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/blendle/zapdriver"
	"github.com/spf13/viper"
)

var (
	logger *zap.Logger
)

func init() {
	verbose := viper.GetBool("VERBOSE")
	if verbose {
		logger, _ = zapdriver.NewDevelopment()
	} else {
		logger, _ = zapdriver.NewProduction()
	}
	defer logger.Sync()
}

// Info passes through an Info level message to zapdriver
func Info(message string, fields ...zapcore.Field) {
	logger.Info(message, fields...)
}

// Debug passes through a Debug level message to zapdriver
func Debug(message string, fields ...zapcore.Field) {
	logger.Debug(message, fields...)
}

// Warn passes through a Warn level message to zapdriver
func Warn(message string, fields ...zapcore.Field) {
	logger.Warn(message, fields...)
}

// Error passes through an Error level message to zapdriver
func Error(message string, fields ...zapcore.Field) {
	logger.Error(message, fields...)
}

// Fatal passes through a Fatal level message to zapdriver
func Fatal(message string, fields ...zapcore.Field) {
	logger.Fatal(message, fields...)
}
