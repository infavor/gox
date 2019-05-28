package cache

import "reflect"

var (
	cacheBytesContainer    = make(map[int]chan ByteCapsule)
	cacheResourceContainer = make(map[reflect.Type]chan interface{})
	cacheBytesSize         = 10
	cacheResourceSize      = 10
)

// SetCacheBytesListSize sets max size of each cache list.
func SetCacheBytesListSize(size int) {
	if size < 0 {
		return
	}
	cacheBytesSize = size
}

// SetCacheResourceListSize sets max size of each cache list.
func SetCacheResourceListSize(size int) {
	if size < 0 {
		return
	}
	cacheResourceSize = size
}

// ByteCapsule is small bytes container
type ByteCapsule struct {
	dynamic bool
	bytes   []byte
	size    int
}

func makeBytes(size int) []byte {
	return make([]byte, size)
}

// ApplyBytes applies specified size of bytes array.
// dynamic bytes apply will not be cached
func ApplyBytes(size int, dynamic bool) ByteCapsule {
	if dynamic || cacheBytesSize <= 0 {
		return ByteCapsule{
			dynamic: true,
			size:    size,
			bytes:   makeBytes(size),
		}
	}
	var bc ByteCapsule
	cha := cacheBytesContainer[size]
	if cha == nil {
		cha = make(chan ByteCapsule, cacheBytesSize)
		cacheBytesContainer[size] = cha
	}
	select {
	case bc = <-cha:
		return bc
	default:
		return ByteCapsule{
			dynamic: dynamic,
			size:    size,
			bytes:   makeBytes(size),
		}
	}
}

// ReCacheBytes caches bytes ByteCapsule
func ReCacheBytes(bc ByteCapsule) {
	if bc.dynamic {
		return
	}
	cha := cacheBytesContainer[bc.size]
	if cha == nil {
		return
	}
	select {
	case cha <- bc:
	default:
		bc.bytes = nil
	}
}

// Bytes returns bytes array of ByteCapsule
func (bc *ByteCapsule) Bytes() []byte {
	return bc.bytes
}

// ApplyResource applies specified type of resource which cached before.
func ApplyResource(p reflect.Type, fallback func() interface{}) interface{} {
	var bc interface{}
	cha := cacheResourceContainer[p]
	if cha == nil {
		cha = make(chan interface{}, cacheResourceSize)
		cacheResourceContainer[p] = cha
	}
	select {
	case bc = <-cha:
		return bc
	default:
		if fallback != nil {
			return fallback()
		}
		return nil
	}
}

// ReCacheResource caches bytes ByteCapsule
func ReCacheResource(res interface{}) {
	if res == nil {
		return
	}
	cha := cacheResourceContainer[reflect.TypeOf(res)]
	if cha == nil {
		return
	}
	select {
	case cha <- res:
	default:
	}
}
