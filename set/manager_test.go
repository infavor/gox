package set_test

import (
	"bytes"
	"fmt"
	"github.com/infavor/gox"
	"github.com/infavor/gox/convert"
	"github.com/infavor/gox/file"
	"github.com/infavor/gox/hash/hashcode"
	"github.com/infavor/gox/logger"
	"github.com/infavor/gox/set"
	"sync"
	"testing"
)

func init() {
	logger.Init(nil)
}

func TestInitHashCode(t *testing.T) {
	for i := 0; i < 10000; i++ {
		key := gox.Md5Sum(convert.IntToStr(i))
		h := hashcode.HashCode(key)
		index := (10000 - 1) & int(h)
		fmt.Println(index)
	}
}

func TestInit(t *testing.T) {
	out, _ := file.CreateFile("D:\\tmp\\godfs\\block")
	out.WriteAt([]byte{0}, 1024)
	out.WriteAt([]byte{1}, 0)
	out.Close()
}

func TestLock(t *testing.T) {
	fmt.Println(111)
	a := make([]chan byte, 10000000)
	fmt.Println(len(a))
	go func() {
		a[0] <- 1
	}()
	<-a[0]

	var g = sync.WaitGroup{}
	g.Add(1)
	g.Wait()
	fmt.Println("done")
}

func TestManagerWrite(t *testing.T) {
	b := 1 << 20
	ss := 8
	manager, err := set.NewFileMap(b, ss, "D:\\tmp\\godfs\\block")
	if err != nil {
		logger.Fatal(err)
	}
	data := []byte{1, 1, 1, 1, 1, 1, 1, 1}
	logger.Info("start")
	for i := 0; i < b; i += 2 {
		if err = manager.Write(i, data); err != nil {
			logger.Fatal(err)
		}
	}
	logger.Info("end")
}

func TestManagerRead(t *testing.T) {
	b := 1 << 20
	ss := 8
	manager, err := set.NewFileMap(b, ss, "D:\\tmp\\godfs\\block")
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("start")
	success := 0
	d := []byte{1, 1, 1, 1, 1, 1, 1, 1}
	for i := 0; i < b; i += 2 {
		data, err := manager.Read(i)
		if err != nil {
			logger.Fatal(err)
		}
		if bytes.Equal(d, data) {
			success++
		}
		//logger.Info(data)
	}
	logger.Info(success)
	logger.Info("end")
}
