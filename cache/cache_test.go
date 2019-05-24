package cache_test

import (
	"fmt"
	"github.com/hetianyi/gox/cache"
	"testing"
	"time"
)

func TestAll(t *testing.T) {
	time.Sleep(time.Second * 1)

	bc := cache.ApplyBytes(102*1, false)
	bc.Bytes()[2] = 11
	cache.ReCacheBytes(bc)
	for i := 0; i < 10; i++ {
		fmt.Println("apply...")
		bc := cache.ApplyBytes(102*1, false)
		fmt.Println(bc.Bytes())
		fmt.Println()
		cache.ReCacheBytes(bc)
	}

	time.Sleep(time.Second * 10)
}
