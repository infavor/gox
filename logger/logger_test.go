package logger_test

import (
	"testing"

	"github.com/infavor/gox/logger"
	"github.com/logrusorgru/aurora"
)

func TestOnlyConsole(t *testing.T) {
	logger.Init(&logger.Config{
		Level:              logger.InfoLevel,
		Write2File:         false,
		AlwaysWriteConsole: true,
	})
	for i := 0; i < 100; i++ {
		logger.Infof("foo is %d", i)
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

func TestBuff(t *testing.T) {
	logger.Init(&logger.Config{
		Level:              logger.InfoLevel,
		Write2File:         true,
		AlwaysWriteConsole: true,
		RollingFileDir:     "D:\\tmp\\logs",
	})
	for i := 0; i < 100; i++ {
		logger.Info(i)
	}
	logger.Sync()
}
