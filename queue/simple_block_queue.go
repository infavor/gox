package queue

type Queue struct {
	size int
	cha  chan interface{}
}

func NewQueue(size int) *Queue {
	return &Queue{
		size: size,
		cha:  make(chan interface{}, size),
	}
}

func (q *Queue) Put(item interface{}) {
	q.cha <- item
}

func (q *Queue) Fetch() interface{} {
	return <-q.cha
}
