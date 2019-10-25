package mapfile_test

import (
	"github.com/hetianyi/gox"
	"github.com/hetianyi/gox/convert"
	"github.com/hetianyi/gox/hash/hashcode"
	"github.com/hetianyi/gox/logger"
	"github.com/hetianyi/gox/mapfile"
	"testing"
)

var (
	manager       *mapfile.FixedSizeFileMap
	ao            *mapfile.AppendFile
	addressBuffer []byte
	slotNum       = 1 << 20
)

func init() {
	ss := 8
	m, err := mapfile.NewFileMap(slotNum, ss, "C:\\k8s\\godfs\\manager")
	if err != nil {
		logger.Fatal(err)
	}
	a, err := mapfile.NewAppendFile(32, 2, "C:\\k8s\\godfs\\aof")
	if err != nil {
		logger.Fatal(err)
	}

	manager = m
	ao = a
	addressBuffer = make([]byte, 8)
}

// --------------
//  write 52973 ms
//  1000000
//  18.878/ms
//  18878/s
// --------------
//  read 23412 ms
//  1000000
//  42.713/ms
//  42713/s
func TestManagedAof(t *testing.T) {

	logger.Info("start write")

	for i := 0; i < 10000; i++ {
		key := gox.Md5Sum(convert.IntToStr(i))
		h := hashcode.HashCode(key)
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
			if err := ao.Write([]byte(key), l); err != nil {
				logger.Fatal(err)
			}
		}
	}

	logger.Info("end write")

	//fmt.Println("xxxxxxxxxxx")

	for i := 0; i < 10000; i++ {
		key := gox.Md5Sum(convert.IntToStr(i))
		h := hashcode.HashCode(key)
		index := (slotNum - 1) & int(h)
		addr, err := manager.Read(index)
		if err != nil {
			logger.Fatal(err)
		}
		l := convert.Bytes2Length(addr)

		_, err = ao.Contains([]byte(key), l)
		if err != nil {
			logger.Fatal(err)
		}
		//fmt.Println(c)
	}

	logger.Info("end read")
}
