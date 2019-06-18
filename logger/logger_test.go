package logger_test

import (
	"github.com/hetianyi/gox/logger"
	"testing"
)

func TestOnlyConsole(t *testing.T) {
	logger.Init(&logger.Config{
		Level:              logger.InfoLevel,
		Write2File:         true,
		AlwaysWriteConsole: true,
	})
	for i := 0; i < 100; i++ {
		logger.Info(i)
	}
}
