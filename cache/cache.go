package cache

import "reflect"

var (
	cacheContainer = make(map[int]chan *ByteCapsule)
	cacheSize      = 10
)

// SetCacheListSize sets max size of each cache list.
func SetCacheListSize(cacheListSize int) {
	if cacheListSize < 0 {
		return
	}
	cacheSize = cacheListSize
}

// ByteCapsule is small bytes container
type ByteCapsule struct {
	dynamic bool
	bytes   []byte
	size    int
}

// ByteCapsule is small bytes container
type ResCapsule struct {
	dynamic bool
	res     interface{}
}

func makeRes(p reflect.Type) interface{} {
	return reflect.New(p).Interface()
}

func makeBytes(size int) []byte {
	return make([]byte, size)
}

// Apply applies specified size of bytes array.
// dynamic bytes apply will not be cached
func ApplyBytes(size int, dynamic bool) *ByteCapsule {
	if dynamic || cacheSize <= 0 {
		return &ByteCapsule{
			dynamic: true,
			size:    size,
			bytes:   makeBytes(size),
		}
	}
	var bc *ByteCapsule
	cha := cacheContainer[size]
	if cha == nil {
		cha = make(chan *ByteCapsule, cacheSize)
		cacheContainer[size] = cha
	}
	select {
	case bc = <-cha:
		return bc
	default:
		return &ByteCapsule{
			dynamic: dynamic,
			size:    size,
			bytes:   makeBytes(size),
		}
	}
}

// ReCache caches bytes ByteCapsule
func ReCacheBytes(bc *ByteCapsule) {
	if bc == nil || bc.dynamic {
		return
	}
	cha := cacheContainer[bc.size]
	if cha == nil {
		return
	}
	select {
	case cha <- bc:
	default:
		bc.bytes = nil
		bc = nil
	}
}

// Bytes returns bytes array of ByteCapsule
func (bc *ByteCapsule) Bytes() []byte {
	return bc.bytes
}

// Apply applies specified size of bytes array.
// dynamic bytes apply will not be cached
func ApplyResource(size int, dynamic bool) *ByteCapsule {

	if dynamic || cacheSize <= 0 {
		return &ByteCapsule{
			dynamic: true,
			size:    size,
			bytes:   makeBytes(size),
		}
	}
	var bc *ByteCapsule
	cha := cacheContainer[size]
	if cha == nil {
		cha = make(chan *ByteCapsule, cacheSize)
		cacheContainer[size] = cha
	}
	select {
	case bc = <-cha:
		return bc
	default:
		return &ByteCapsule{
			dynamic: dynamic,
			size:    size,
			bytes:   makeBytes(size),
		}
	}
}

// ReCache caches bytes ByteCapsule
func ReCacheRes(res interface{}) {
	if bc == nil || bc.dynamic {
		return
	}
	cha := cacheContainer[bc.size]
	if cha == nil {
		return
	}
	select {
	case cha <- bc:
	default:
		bc.bytes = nil
		bc = nil
	}
}
