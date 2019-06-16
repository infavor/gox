package hashcode_test

import (
	"fmt"
	"github.com/hetianyi/gox/hash/hashcode"
	"testing"
)

func TestHashCode(t *testing.T) {
	fmt.Println(hashcode.HashCode("kjasdj阿斯达斯的1232ASDad-==：‘’"))
	fmt.Println(hashcode.HashCode(1))
	fmt.Println(hashcode.HashCode(int64(3000000000000000)))
	fmt.Println(hashcode.HashCode(float32(3000000000000000)))
}
