package cache

import "fmt"

var (
	cacheContainer = make(map[int]chan *ByteCapsule)
)

type ByteCapsule struct {
	dynamic bool
	bytes   []byte
	size    int
}

func makeBuffer(size int) []byte {
	fmt.Println("make buffer...")
	return make([]byte, size)
}

// Apply applies specified size of bytes array.
// dynamic bytes apply will not be cached
func Apply(size int, dynamic bool) *ByteCapsule {
	if dynamic {
		return &ByteCapsule{
			dynamic: true,
			size:    size,
			bytes:   makeBuffer(size),
		}
	}
	var bc *ByteCapsule
	cha := cacheContainer[size]
	if cha == nil {
		cha = make(chan *ByteCapsule, 10)
		cacheContainer[size] = cha
	}
	select {
	case bc = <-cha:
		return bc
	default:
		return &ByteCapsule{
			dynamic: dynamic,
			size:    size,
			bytes:   makeBuffer(size),
		}
	}
}

func Recache(bc *ByteCapsule) {
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

//
func (bc *ByteCapsule) Bytes() []byte {
	return bc.bytes
}
