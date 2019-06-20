package main

import (
	"github.com/hetianyi/gox/logger"
	"github.com/hetianyi/gox/uuid"
	"time"
)

func main() {
	logger.Init(&logger.Config{
		Level:              logger.InfoLevel,
		RollingPolicy:      []int{logger.HOUR, logger.MB64},
		Write2File:         true,
		AlwaysWriteConsole: false,
		RollingFileDir:     "D:\\Hetianyi\\svn\\gox\\logger\\tes",
		RollingFileName:    "godfs",
	})

	//go changeLength()
	for true {
		logger.Info(uuid.UUID())
	}

}

func changeLength() {
	var le int64 = 2 << 25
	for true {
		time.Sleep(time.Second * 16)
		logger.FakeWriteLen(le)
		le = le << 1
	}
}
