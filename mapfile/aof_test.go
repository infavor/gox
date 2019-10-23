package mapfile_test

import (
	"github.com/hetianyi/gox/file"
	"github.com/hetianyi/gox/logger"
	"testing"
)

func init() {
	logger.Init(nil)
}

func TestInitAOF(t *testing.T) {
	out, _ := file.CreateFile("D:\\tmp\\godfs\\aof")

}
