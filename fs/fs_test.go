package fs

import (
	"fmt"
	"github.com/hetianyi/gox"
	"testing"
)

func Test1(t *testing.T) {
	s := "asd123"
	var h int32 = 0
	if len(s) > 0 {
		for _, r := range s {
			h = 31*h + r
		}
	}
	fmt.Println("HasCode:", h)
}
func Test2(t *testing.T) {
	fmt.Println(gox.HashCode(""))
	fmt.Println(gox.HashCode("1"))
	fmt.Println(gox.HashCode("123"))
	fmt.Println(gox.HashCode("abc"))
	fmt.Println(gox.HashCode("abc123"))
	fmt.Println(gox.HashCode("kasdasfoas卡死扩大宽松的阿萨德2312312323"))
}

func Test3(t *testing.T) {
	m := make(map[string]interface{}, 2)
	fmt.Println(len(m))
	s := "asd123"
	hash := gox.HashCode(s)
	i := int((2 - 1) & int(hash))
	fmt.Println("index=", i)
}
