package queue_test

import (
	"fmt"
	"github.com/infavor/gox/queue"
	"testing"
	"time"
)

func TestNewQueue(t *testing.T) {
	q := queue.NewQueue(2)

	go producer(q)
	go consumer(q)

	time.Sleep(time.Second * 20)
}
func ExampleNewQueue() {
	q := queue.NewQueue(2)

	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println("添加：", i)
			q.Put(i)
			fmt.Println("添加结束：", i)
		}
	}()
	go func() {
		for true {
			fmt.Println("得到：", q.Fetch())
			time.Sleep(time.Second)
		}
	}()
	time.Sleep(time.Second * 20)
}

func TestNewQueueExample(t *testing.T) {
	q := queue.NewQueue(2)

	go producer(q)
	go consumer(q)

	time.Sleep(time.Second * 20)
}

func producer(q *queue.Queue) {
	for i := 0; i < 10; i++ {
		fmt.Println("添加：", i)
		q.Put(i)
		fmt.Println("添加结束：", i)
	}
}

func consumer(q *queue.Queue) {
	for true {
		fmt.Println("得到：", q.Fetch())
		time.Sleep(time.Second)
	}
}

func TestQueue_Put(t *testing.T) {
	q := queue.NewNoneBlockQueue(2)
	fmt.Println(q.Put(1))
	fmt.Println(q.Put(1))
	fmt.Println(q.Put(1))
	fmt.Println(q.Fetch())
	fmt.Println(q.Fetch())
	fmt.Println(q.Fetch())
	fmt.Println(q.Fetch())
	fmt.Println(q.Put(1))
	fmt.Println(q.Fetch())
}
