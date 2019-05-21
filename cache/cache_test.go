package cache_test

import (
	"../cache"
	"fmt"
	"testing"
	"time"
)

func TestAll(t *testing.T) {
	time.Sleep(time.Second * 1)

	bc := cache.Apply(102*1, false)
	bc.Bytes()[2] = 11
	cache.Recache(bc)
	for i := 0; i < 10; i++ {
		fmt.Println("apply...")
		bc := cache.Apply(102*1, false)
		fmt.Println(bc.Bytes())
		fmt.Println()
		cache.Recache(bc)
	}

	time.Sleep(time.Second * 10)
}
