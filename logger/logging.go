package logger

import "go.uber.org/zap"

var logger *zap.Logger

// see https://github.com/betterstack-community/go-logging
func Get() *zap.Logger {
	logger = zap.Must(zap.NewDevelopment())
	return logger
}
