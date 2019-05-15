package mem

type ByteCapsule struct {
	dynamic bool
	Bytes   []byte
}

func makeBuffer(size int) []byte {
	return make([]byte, size)
}

// Apply applies specified size of bytes array.
// dynamic bytes apply will not be cached
func Apply(size int, dynamic bool) *ByteCapsule {
	if dynamic {
		return &ByteCapsule{
			dynamic: true,
			Bytes:   makeBuffer(size),
		}
	}

}
