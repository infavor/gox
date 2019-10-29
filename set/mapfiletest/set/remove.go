package main

import (
	"fmt"
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
		slotNum  int = 1 << 25
		slotSize int
		caseSize int
	)

	caseSize = 1000000

	slotSize = 32 * 3
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

	logger.Info("start remove")

	for i := 0; i < caseSize; i++ {
		key := gox.Md5Sum(convert.IntToStr(i))
		key = key + key + key
		c, err := set.Remove([]byte(key))
		if err != nil {
			logger.Fatal(err)
		}
		if !c {
			logger.Error("delete failed:", c)
		}
	}

	logger.Info("end remove")
	total := 0
	for _, v := range manager.SlotSnapshot() {
		total += int(v)
	}

	fmt.Println("empty slots: ", slotNum-total)
	fmt.Println("slots usage: ", convert.Float64ToStr(float64(total) / float64(slotNum) * 100)[0:6]+"%")
}
