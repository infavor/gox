package gox_test

import (
	"github.com/hetianyi/gox"
	"github.com/hetianyi/gox/logger"
	"testing"
)

func init() {
	logger.Init(nil)
}

func TestNetwork(t *testing.T) {
	logger.Info(gox.GetMyAddress("vEthernet", "192.168.0"))
}
