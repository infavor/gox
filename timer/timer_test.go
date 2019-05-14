package timer_test

import (
	"../timer"
	"fmt"
	"testing"
	"time"
)

/*func Test1(t *testing.T) {
	fmt.Println(time.Now().Second())
	timer.Start(5000, 1000, 0, func() {
		time.Sleep(time.Second * 3)
		fmt.Println("xxx", time.Now().Second())
	})
	c := make(chan int)
	<- c
}*/

func Test2(t *testing.T) {
	fmt.Println(time.Now().Second())
	timer.Start(4000, 0, 1000, func() {
		fmt.Println("xxx", time.Now().Second())
	})
	c := make(chan int)
	<-c
}
