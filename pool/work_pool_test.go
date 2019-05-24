package pool_test

import (
	"fmt"
	"github.com/hetianyi/gox/pool"
	"testing"
)

func Test1(t *testing.T) {
	p := pool.New(97, 100)
	var i = 0
	for {
		i++
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
	fmt.Println("execute task ", taskId)
}
