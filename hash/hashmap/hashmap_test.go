package hashmap_test

import (
	"fmt"
	"github.com/hetianyi/gox"
	"github.com/hetianyi/gox/convert"
	"github.com/hetianyi/gox/hash/hashmap"
	"github.com/hetianyi/gox/logger"
	"github.com/sirupsen/logrus"
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
	time.Sleep(time.Second * 5)
	logrus.Info("开始添加....")
	for i := 0; i < 10000000; i++ {
		m.Put(convert.IntToStr(i), i)
	}
	logrus.Info("结束添加....")
	time.Sleep(time.Second * 10)

	logrus.Info("测试查找1....")
	logrus.Info("1=", m.Get("1"))
	logrus.Info("测试查找2....")
	logrus.Info("2=", m.Get("2"))
	logrus.Info("结束查找....")
	time.Sleep(time.Second * 10)
}

func TestNewMap2(t *testing.T) {
	m := hashmap.NewMap()
	m.Put("123", "123")
	m.Put("123", "123")
}
