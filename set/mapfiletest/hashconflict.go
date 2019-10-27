package main

import (
	"fmt"
	"github.com/hetianyi/gox"
	"github.com/hetianyi/gox/convert"
	"github.com/hetianyi/gox/hash/hashcode"
	"github.com/hetianyi/gox/logger"
	"os"
)

func main() {
	var (
		slotNum  int // 16777216
		caseSize int
		hashs    = make(map[int]int)
		slotMap  []byte
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
	slotMap = make([]byte, slotNum)

	_caseSize, err := convert.StrToInt(cs)
	if err != nil {
		logger.Fatal(err)
	}
	caseSize = _caseSize

	logger.Info("start")

	for i := 0; i < caseSize; i++ {
		key := gox.Md5Sum(convert.IntToStr(i))
		h := hashcode.HashCode(key)
		h ^= h >> 16
		index := (slotNum - 1) & int(h)
		slotMap[index] = 1
		hashs[index] = hashs[index] + 1
	}

	logger.Info("end")

	total := 0
	for _, v := range slotMap {
		total += int(v)
	}

	fmt.Println("empty slots: ", slotNum-total)
	fmt.Println("slots usage: ", convert.Float64ToStr(float64(total) / float64(slotNum) * 100)[0:6]+"%")
	fmt.Println("unique hash: ", len(hashs))
	fmt.Println("collision percent: ", convert.Float64ToStr(float64(caseSize-len(hashs))*100/float64(caseSize))+"%")
}
