package main

import "github.com/hetianyi/gox/logger"

func main() {
	logger.Init(&logger.Config{
		Level:              logger.InfoLevel,
		Write2File:         false,
		AlwaysWriteConsole: true,
		RollingFileDir:     "/tmp",
		RollingFileName:    "FUCK",
	})
	logger.Info("xxxxxxxxxxxx\n123123123")
}
