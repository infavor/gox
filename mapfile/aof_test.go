package mapfile_test

import (
	"github.com/hetianyi/gox/file"
	"github.com/hetianyi/gox/logger"
	"github.com/hetianyi/gox/mapfile"
	"testing"
)

func init() {
	logger.Init(nil)
}

func TestInitAOF(t *testing.T) {
	a, err := mapfile.NewAppendFile(5, 1, "D:\\tmp\\godfs\\aof")
	if err != nil {
		logger.Fatal(err)
	}
	a.Write("11111")
}
