package set_test

import (
	"fmt"
	"github.com/hetianyi/gox"
	"github.com/hetianyi/gox/logger"
	"github.com/hetianyi/gox/set"
	"github.com/hetianyi/gox/timer"
	"os"
	"sync"
	"testing"
	"time"
)

func init() {
	logger.Init(nil)
}

func TestInitAOF(t *testing.T) {
	a, err := set.NewAppendFile(5, 1, "D:\\tmp\\godfs\\aof")
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
	a, err := set.NewAppendFile(5, 1, "D:\\tmp\\godfs\\aof")
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
	a, err := set.NewAppendFile(5, 1, "D:\\tmp\\godfs\\aof")
	if err != nil {
		logger.Fatal(err)
	}
	var addr int64 = 0
	x, _, err := a.Contains([]byte("11111"), addr)
	if err != nil {
		logger.Fatal(err)
	}
	fmt.Println(x)
}

func TestInitAOFContains1(t *testing.T) {
	a, err := set.NewAppendFile(5, 1, "D:\\tmp\\godfs\\aof")
	if err != nil {
		logger.Fatal(err)
	}
	var addr int64 = 0

	logger.Info("start")
	for i := 0; i < 1000; i++ {
		x, _, err := a.Contains([]byte("11111"), addr)
		if err != nil {
			logger.Fatal(err)
		}
		fmt.Println(x)
	}
	logger.Info("end")
	// 1000/22ms = 41/ms
}

func TestInitAOFContains2(t *testing.T) {
	a, err := set.NewAppendFile(5, 1, "D:\\tmp\\godfs\\aof")
	if err != nil {
		logger.Fatal(err)
	}
	var addr int64 = 0

	logger.Info("start")

	lock := new(sync.Mutex)
	val := 0
	inc := func() {
		lock.Lock()
		defer lock.Unlock()
		val++
		if val == 50000 {
			logger.Info("finish")
			os.Exit(0)
		}
	}

	q := func() {
		for i := 0; i < 1000; i++ {
			_, _, err := a.Contains([]byte("11111"), addr)
			if err != nil {
				logger.Fatal(err)
			}
			inc()
		}
	}

	for i := 0; i < 50; i++ {
		go q()
	}

	logger.Info("end")

	timer.Start(0, time.Second, 0, func(t *timer.Timer) {
		fmt.Println(val)
	})
	gox.BlockTest()
	// 1000/22ms = 41/ms
}
