package pool_test

import (
	"../pool"
	"fmt"
	"testing"
	"time"
)

func Test1(t *testing.T) {
	p := pool.New(97, 100)
	for i := 0; i < 100; i++ {
		tmp := i
		err := p.Push(func() {
			testTask(tmp)
		})
		if err != nil {
			fmt.Println("Err:", err)
		}
	}
	c := make(chan int)
	<-c
}

func testTask(taskId int) {
	time.Sleep(time.Second * 1)
	fmt.Println("execute task ", taskId)
}
