package gox_test

import (
	"github.com/hetianyi/gox"
	"github.com/hetianyi/gox/logger"
	"github.com/sirupsen/logrus"
	"testing"
)

func init() {
	logger.Init(nil)
}

func TestNetwork(t *testing.T) {
	logrus.Info(gox.GetMyAddress("vEthernet", "192.168.0"))
}
