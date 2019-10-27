package main

import (
	"github.com/hetianyi/gox"
	"github.com/hetianyi/gox/convert"
	"github.com/hetianyi/gox/logger"
	"github.com/hetianyi/gox/set"
)

func init() {
	logger.Init(nil)
}

func main() {
	var (
		manager  *set.FixedSizeFileMap
		ao       *set.AppendFile
		slotNum  int
		slotSize int
		caseSize int
	)

	slotNum = 1 << 20

	caseSize = 500000

	slotSize = 32
	m, err := set.NewFileMap(slotNum, 8, "index")
	if err != nil {
		logger.Fatal(err)
	}
	a, err := set.NewAppendFile(slotSize, 2, "aof")
	if err != nil {
		logger.Fatal(err)
	}

	manager = m
	ao = a
	set := set.NewDataSet(manager, ao)

	logger.Info("start reading")

	for i := 0; i < caseSize; i++ {
		key := gox.Md5Sum(convert.IntToStr(i))
		c, err := set.Contains([]byte(key))
		if err != nil {
			logger.Fatal(err)
		}
		if !c {
			logger.Fatal("check failed:", i)
		}
	}

	logger.Info("end reading")
}
