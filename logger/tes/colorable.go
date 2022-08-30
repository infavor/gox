package main

import (
	"github.com/infavor/gox/logger"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
)

func main() {

	logrus.SetFormatter(&logger.SimpleTextFormatter{})
	logrus.SetOutput(colorable.NewColorableStdout())

	logger.Info("succeeded")
	logrus.Info("succeeded")
	logrus.Warn("not correct")
	logrus.Error("something error")
	logrus.Fatal("panic")
}
