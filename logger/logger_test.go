package logger_test

import (
	"github.com/hetianyi/gox/logger"
	"github.com/logrusorgru/aurora"
	"testing"
)

func TestOnlyConsole(t *testing.T) {
	logger.Init(&logger.Config{
		Level:              logger.InfoLevel,
		Write2File:         false,
		AlwaysWriteConsole: true,
	})
	for i := 0; i < 100; i++ {
		logger.Info(i)
	}
}

func TestSimpleTextFormatter_Format(t *testing.T) {
	logger.Init(&logger.Config{
		Level:     logger.InfoLevel,
		Formatter: &logger.SimpleTextFormatter{},
	})
	logger.Info("Hello")
}

func TestNoneTextFormatter_Format(t *testing.T) {
	logger.Init(&logger.Config{
		Level:     logger.InfoLevel,
		Formatter: &logger.NoneTextFormatter{},
	})
	logger.Info("Hello")
}

func TestShortTextFormatter_Format(t *testing.T) {
	logger.Init(&logger.Config{
		Level:     logger.InfoLevel,
		Formatter: &logger.ShortTextFormatter{},
	})
	logger.Info("Hello")
}

func TestPrintColor(t *testing.T) {
	logger.PrintColor([]byte(aurora.Cyan("Hello").String()))
}
