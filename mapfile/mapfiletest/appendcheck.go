package main

import (
	"github.com/hetianyi/gox"
	"github.com/hetianyi/gox/convert"
	"github.com/hetianyi/gox/hash/hashcode"
	"github.com/hetianyi/gox/logger"
	"github.com/hetianyi/gox/mapfile"
	"os"
)

func main() {
	var (
		manager       *mapfile.FixedSizeFileMap
		ao            *mapfile.AppendFile
		addressBuffer []byte
		slotNum       int
		slotSize      int
		caseSize      int
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

	m, err := mapfile.NewFileMap(slotNum, 8, "index")
	if err != nil {
		logger.Fatal(err)
	}
	a, err := mapfile.NewAppendFile(slotSize, 2, "aof")
	if err != nil {
		logger.Fatal(err)
	}

	manager = m
	ao = a
	addressBuffer = make([]byte, 8)

	logger.Info("start writing")

	for i := 0; i < caseSize; i++ {
		key := gox.Md5Sum(convert.IntToStr(i))
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

	logger.Info("end writing")
	logger.Info("start reading")

	for i := 0; i < caseSize; i++ {
		key := gox.Md5Sum(convert.IntToStr(i))
		h := hashcode.HashCode(key)
		index := (slotNum - 1) & int(h)
		addr, err := manager.Read(index)
		if err != nil {
			logger.Fatal(err)
		}
		l := convert.Bytes2Length(addr)

		_, _, err = ao.Contains([]byte(key), l)
		if err != nil {
			logger.Fatal(err)
		}
	}

	logger.Info("end reading")

}
