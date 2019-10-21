package binlog_test

import (
	"fmt"
	"github.com/hetianyi/gox/file"
	"sync"
	"testing"
)

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
