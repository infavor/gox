package mapfile_test

import (
	"fmt"
	"github.com/hetianyi/gox/logger"
	"github.com/hetianyi/gox/mapfile"
	"sync"
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

func TestInitAOFContains(t *testing.T) {
	a, err := mapfile.NewAppendFile(5, 1, "D:\\tmp\\godfs\\aof")
	if err != nil {
		logger.Fatal(err)
	}
	var addr int64 = 0
	x, err := a.Contains([]byte("11111"), addr)
	if err != nil {
		logger.Fatal(err)
	}
	fmt.Println(x)
}

func TestInitAOFContains1(t *testing.T) {
	a, err := mapfile.NewAppendFile(5, 1, "D:\\tmp\\godfs\\aof")
	if err != nil {
		logger.Fatal(err)
	}
	var addr int64 = 0

	logger.Info("start")
	for i := 0; i < 1000; i++ {
		x, err := a.Contains([]byte("11111"), addr)
		if err != nil {
			logger.Fatal(err)
		}
		fmt.Println(x)
	}
	logger.Info("end")
	// 1000/22ms = 41/ms
}

func TestInitAOFContains2(t *testing.T) {
	a, err := mapfile.NewAppendFile(5, 1, "D:\\tmp\\godfs\\aof")
	if err != nil {
		logger.Fatal(err)
	}
	var addr int64 = 0

	logger.Info("start")

	q := func() {
		for i := 0; i < 1000; i++ {
			_, err := a.Contains([]byte("11111"), addr)
			if err != nil {
				logger.Fatal(err)
			}
		}
	}

	var g = sync.WaitGroup{}
	g.Add(5)

	for i := 0; i < 50; i++ {
		go q()
	}

	g.Wait()

	logger.Info("end")
	<-make(chan int)
	// 1000/22ms = 41/ms
}
