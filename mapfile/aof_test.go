package mapfile_test

import (
	"fmt"
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
	addr, err := a.ApplyAddress()
	if err != nil {
		logger.Fatal(err)
	}
	if err = a.Write([]byte("11111"), addr); err != nil {
		logger.Fatal(err)
	}
	fmt.Println(1)
}

func TestInitAOFSameSlot(t *testing.T) {
	a, err := mapfile.NewAppendFile(5, 1, "D:\\tmp\\godfs\\aof")
	if err != nil {
		logger.Fatal(err)
	}
	var addr int64 = 0
	if err = a.Write([]byte("22222"), addr); err != nil {
		logger.Fatal(err)
	}
	fmt.Println(1)
}
