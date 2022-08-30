package main

import (
	"github.com/infavor/gox"
	"github.com/infavor/gox/convert"
	"github.com/infavor/gox/hash/hashcode"
	"github.com/infavor/gox/logger"
	"github.com/infavor/gox/set"
	"os"
	"sync"
)

var (
	manager       *set.FixedSizeFileMap
	ao            *set.AppendFile
	addressBuffer []byte
	slotNum       int
	slotSize      int
	caseSize      int
)

func main() {

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
	thread := 10
	addressBuffer = make([]byte, 8)

	g := sync.WaitGroup{}
	g.Add(thread)

	saveRange := func(start, end int) {
		for k := start; k < end; k++ {
			save(k)
		}
		g.Done()
	}

	logger.Info("start writing")
	step := caseSize / thread
	for i := 0; i < thread; i++ {
		start := i * step
		end := (i + 1) * step
		if i == thread-1 {
			end = caseSize
		}
		go saveRange(start, end)
	}
	g.Wait()
	logger.Info("end writing")
}

func save(val int) {
	key := gox.Md5Sum(convert.IntToStr(val))
	h := hashcode.HashCode(key)
	h ^= h >> 16
	index := (slotNum - 1) & int(h)
	addr, err := manager.Read(index)
	if err != nil {
		logger.Fatal(err)
	}
	var l int64 = 0
	if addr != nil {
		l = convert.Bytes2Length(addr)
	}

	if l == 0 {
		addr, err := ao.ApplyAddress()
		if err != nil {
			logger.Fatal(err)
		}
		if err := manager.Write(index, convert.Length2Bytes(addr, addressBuffer)); err != nil {
			logger.Fatal(err)
		}
		if err := ao.Write([]byte(key), addr); err != nil {
			logger.Fatal(err)
		}
	} else {
		x, _, err := ao.Contains([]byte(key), l)
		if err != nil {
			logger.Fatal(err)
		}
		if !x {
			if err := ao.Write([]byte(key), l); err != nil {
				logger.Fatal(err)
			}
		}
	}
}
