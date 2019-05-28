package cache_test

import (
	"github.com/hetianyi/gox/cache"
	"testing"
	"time"
)

func TestAll(t *testing.T) {
	time.Sleep(time.Second * 1)

	for i := 0; i < 10000; i++ {
		bc := cache.ApplyBytes(1024*10240, false)
		cache.ReCacheBytes(bc)
	}

	time.Sleep(time.Second * 10)
}
