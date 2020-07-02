package queue

type NoneBlockQueue struct {
	size int
	cha  chan interface{}
}

func NewNoneBlockQueue(size int) *NoneBlockQueue {
	return &NoneBlockQueue{
		size: size,
		cha:  make(chan interface{}, size),
	}
}

func (q *NoneBlockQueue) Put(item interface{}) bool {
	select {
	case q.cha <- item:
		return true
	default:
		return false
	}
}

func (q *NoneBlockQueue) Fetch() (interface{}, bool) {
	var ret interface{}
	select {
	case ret = <-q.cha:
		return ret, true
	default:
		return nil, false
	}
}
