package main

import (
	"github.com/hetianyi/gox"
	"github.com/hetianyi/gox/convert"
	"github.com/hetianyi/gox/logger"
	"github.com/hetianyi/gox/set"
	"sync"
)

func init() {
	logger.Init(nil)
}

func main() {
	var (
		manager  *set.FixedSizeFileMap
		ao       *set.AppendFile
		slotNum  int = 1 << 24
		slotSize int
		caseSize int
	)

	caseSize = 10000000

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

	set := set.NewDataSet(manager, ao)

	g := sync.WaitGroup{}
	g.Add(thread)

	saveRange := func(start, end int) {
		logger.Info("writing from ", start, " to ", end)
		for k := start; k < end; k++ {
			key := gox.Md5Sum(convert.IntToStr(k))
			if err := set.Add([]byte(key)); err != nil {
				logger.Fatal(err)
			}
			c, err := set.Contains([]byte(key))
			if err != nil {
				logger.Fatal(err)
			}
			if !c {
				logger.Fatal("cannot get ", k)
			}
			incr()
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
	logger.Info(total)
}

var lock = new(sync.Mutex)
var total = 0

func incr() {
	lock.Lock()
	defer lock.Unlock()
	total++
}
