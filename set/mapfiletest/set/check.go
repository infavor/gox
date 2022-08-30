package main

import (
	"fmt"
	"github.com/infavor/gox"
	"github.com/infavor/gox/convert"
	"github.com/infavor/gox/hash/hashcode"
	"github.com/infavor/gox/logger"
	"github.com/infavor/gox/set"
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
	loss := 0

	logger.Info("start reading")

	for i := 0; i < caseSize; i++ {
		key := gox.Md5Sum(convert.IntToStr(i))
		key = key + key + key
		c, err := set.Contains([]byte(key))
		if err != nil {
			logger.Fatal(err)
		}
		if !c {
			loss++
			//logger.Fatal("check failed:", i)
		}
	}

	logger.Info("end reading")
	total := 0
	for _, v := range manager.SlotSnapshot() {
		total += int(v)
	}

	fmt.Println("empty slots : ", slotNum-total)
	fmt.Println("check failed: ", loss)
	fmt.Println("slots usage : ", convert.Float64ToStr(float64(total) / float64(slotNum) * 100)[0:6]+"%")
}

func getIndex(manager *set.FixedSizeFileMap, data []byte) int {
	key := gox.Md5Sum(string(data))
	h := hashcode.HashCode(key)
	h = h ^ (h >> 16)
	return (manager.SlotNum() - 1) & int(h)
}
