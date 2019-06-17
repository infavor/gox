package hashcode_test

import (
	"fmt"
	"github.com/hetianyi/gox/hash/hashcode"
	"reflect"
	"testing"
)

func TestHashCode(t *testing.T) {
	fmt.Println(hashcode.HashCode("kjasdj阿斯达斯的1232ASDad-==：‘’"))
	fmt.Println(hashcode.HashCode(1))
	fmt.Println(hashcode.HashCode(int64(3000000000000000)))
	fmt.Println(hashcode.HashCode(float32(3000000000000000)))
}

func TestHashCode1(t *testing.T) {
	var a byte = 111
	fmt.Println(int32(a))

	var b int8 = -122
	fmt.Println(int32(b))

	var c int16 = -1450
	fmt.Println(int32(c))

	var d int32 = -1450678
	fmt.Println(int32(d))

	var e rune = 'e'
	fmt.Println(int32(e))

	var f uint8 = 255
	fmt.Println(int32(f))

	var g uint = 21334234
	fmt.Println(int32(g))

	var h uint32 = 3133434000
	fmt.Println(int32(h))

	var i int64 = -9133433341231224323
	fmt.Println((int32)(i ^ (i >> 32)))

}

type XR1 struct {
}

func TestHashCode2(t *testing.T) {
	a := make([]XR1, 2)
	fmt.Println(reflect.TypeOf(a))
}

func TestHashCode3(t *testing.T) {
	a := make(map[interface{}]int, 2)
	item := &XR1{}
	fmt.Println(item)
	a[item] = 1
	fmt.Println(a[item])

}
