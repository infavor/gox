package main

import (
	"encoding/json"
	"fmt"
	"github.com/infavor/gox"
	"github.com/infavor/gox/convert"
	"github.com/infavor/gox/hash/hashcode"
	"github.com/infavor/gox/logger"
	"github.com/infavor/gox/set"
	"os"
)

func main() {
	var (
		manager  *set.FixedSizeFileMap
		ao       *set.AppendFile
		slotNum  int
		slotSize int
		caseSize int
		depthMap = make(map[int]int)
	)

	logger.Init(nil)

	if len(os.Args) < 3 {
		logger.Fatal("Usage: ./<app> <slotNum> <caseSize>")
	}
	sn := os.Args[1]
	cs := os.Args[2]

	_slotNum, err := convert.StrToInt(sn)
	if err != nil {
		logger.Fatal(err)
	}
	slotNum = _slotNum

	_caseSize, err := convert.StrToInt(cs)
	if err != nil {
		logger.Fatal(err)
	}
	caseSize = _caseSize

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

	logger.Info("start reading")

	for i := 0; i < caseSize; i++ {
		key := gox.Md5Sum(convert.IntToStr(i))
		h := hashcode.HashCode(key)
		h = h ^ (h >> 16)
		index := (slotNum - 1) & int(h)
		addr, err := manager.Read(index)
		if err != nil {
			logger.Fatal(err)
		}
		l := convert.Bytes2Length(addr)

		_, dep, err := ao.Contains([]byte(key), l)
		if err != nil {
			logger.Fatal(err)
		}
		if dep > 5 {
			depthMap[dep] = depthMap[dep] + 1
		}
	}

	logger.Info("end reading")

	snapshot := manager.SlotSnapshot()
	total := 0
	for _, v := range snapshot {
		total += int(v)
	}

	fmt.Println("empty slots: ", slotNum-total)
	fmt.Println("slots usage: ", convert.Float64ToStr(float64(total) / float64(slotNum) * 100)[0:6]+"%")

	ret, _ := json.MarshalIndent(depthMap, "", "    ")
	fmt.Println(string(ret))
}
