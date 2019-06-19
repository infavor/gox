package hashmap_test

import (
	"fmt"
	"github.com/hetianyi/gox"
	"github.com/hetianyi/gox/convert"
	"github.com/hetianyi/gox/hash/hashmap"
	"github.com/hetianyi/gox/logger"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func init() {
	logger.Init(nil)
}

func TestNewMap(t *testing.T) {
	m := hashmap.NewMap()
	m.Put("123", "asd")
	m.Put("456", "ghj")

	fmt.Println(m.Get("123"))
	fmt.Println(m.Get("456"))
	fmt.Println(1 << 30)
	fmt.Println(2 << 1)

	n := 9 - 1
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	fmt.Println(gox.TValue(n < 0, 1, gox.TValue(n >= 1073741824, 1073741824, n+1)))

}

func TestNewMap1(t *testing.T) {
	m := hashmap.NewMap()
	//time.Sleep(time.Second * 5)
	logrus.Info("开始添加....")
	for i := 0; i < 10000000; i++ {
		m.Put(convert.IntToStr(i), i)
	}
	logrus.Info("结束添加....")
	//time.Sleep(time.Second * 10)

	for i := 0; i < 10000000; i += 99999 {
		start := time.Now()
		logrus.Info(i, "=", m.Get(convert.IntToStr(i)))
		logrus.Info("time:", time.Now().UnixNano()-start.UnixNano())
	}

	fmt.Println("remove 1=", m.Remove("1"))
	fmt.Println("remove 1=", m.Remove("1"))
	fmt.Println("remove 2=", m.Remove("2"))

	logrus.Info("结束查找....")
	time.Sleep(time.Second * 10)
}

func TestBuildInMap(t *testing.T) {
	m := make(map[string]int)
	//time.Sleep(time.Second * 5)
	logrus.Info("开始添加....")
	for i := 0; i < 10000000; i++ {
		m[convert.IntToStr(i)] = i
	}
	logrus.Info("结束添加....")
	//time.Sleep(time.Second * 10)

	for i := 0; i < 10000000; i += 99999 {
		start := time.Now()
		logrus.Info(i, "=", m[convert.IntToStr(i)])
		logrus.Info("time:", time.Now().UnixNano()-start.UnixNano())
	}
	delete(m, "1")
	delete(m, "1")
	delete(m, "2")

	logrus.Info("结束查找....")
	time.Sleep(time.Second * 10)
}

func TestNewMap2(t *testing.T) {
	var a interface{} = 1
	var b interface{} = "x"
	fmt.Println(a == b)
}

func TestNewMap3(t *testing.T) {
	m := hashmap.NewMap()
	for i := 0; i < 150; i++ {
		m.Put(convert.IntToStr(i), i)
	}
}

func TestNewMap5(t *testing.T) {
	m := hashmap.NewMap()
	m.Put("a", "1")
	logrus.Info("a=", m.Get("a"))
	m.Put("a", nil)
	logrus.Info("a=", m.Get("a"))

	var a = 'e'
	m.Put(a, 123)
	logrus.Info("a=", m.Get(a))
}

func myMapRemove() {
	m := hashmap.New(5, 1.0)
	//time.Sleep(time.Second * 5)
	logrus.Info("开始添加....")
	for i := 0; i < 7000000; i++ {
		m.Put(convert.IntToStr(i), i)
	}
	logrus.Info("结束添加....")

	time.Sleep(time.Second * 10)
	logrus.Info("开始remove....")
	for i := 0; i < 7000000; i++ {
		m.Remove(convert.IntToStr(i))
	}
	logrus.Info("结束remove....")
}

func TestNewMapRemove(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	myMapRemove()
	runtime.GC()
	time.Sleep(time.Second * 10000)
}

func buildInMapRemove() {
	m := make(map[string]int)
	//time.Sleep(time.Second * 5)
	logrus.Info("开始添加....")
	for i := 0; i < 7000000; i++ {
		m[convert.IntToStr(i)] = i
	}
	logrus.Info("结束添加....")
	//time.Sleep(time.Second * 10)

	time.Sleep(time.Second * 10)
	logrus.Info("开始remove....")
	for i := 0; i < 7000000; i++ {
		delete(m, convert.IntToStr(i))
	}
	logrus.Info("结束remove....")
}

func TestBuildInMapRemove(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	buildInMapRemove()
	runtime.GC()
	time.Sleep(time.Second * 10000)
}

func TestTableSizeFor(t *testing.T) {
	fmt.Println(tableSizeFor(123))
	a := make(map[string]int)
	a["asd"] = 2
	fmt.Println(a)
}

func tableSizeFor(cap int) int {
	n := cap - 1
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	return gox.TValue(n < 0, 1, gox.TValue(n >= 1073741824, 1073741824, n+1).(int)).(int)
}

func TestInit(t *testing.T) {
	a := 1 | 7
	fmt.Println(a | 3)
}

func Test111(t *testing.T) {
	fmt.Println(filepath.Base("D:/Hetianyi/svn/gox/logger/tes\\godfs-2019061911-part15.log"))
	fmt.Println(filepath.Dir("godfs-2019061911-part15.log"))
}
