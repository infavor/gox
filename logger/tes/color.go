package main

import (
	"github.com/hetianyi/gox/logger"
)

func main() {
	logger.Init(&logger.Config{
		Level:              logger.TraceLevel,
		RollingPolicy:      []int{logger.HOUR, logger.MB64},
		Write2File:         false,
		AlwaysWriteConsole: true,
	})

	logger.Trace("Hello World! 你好，世界！")
	logger.Debug("Hello World! 你好，世界！")
	logger.Info("Hello World! 你好，世界！")
	logger.Warn("Hello World! 你好，世界！")
	logger.Error("Hello World! 你好，世界！")
	logger.Fatal("Hello World! 你好，世界！")
	// logger.Panic("Hello World! 你好，世界！")
}
