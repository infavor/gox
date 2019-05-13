package timer_test

import (
	"../timer"
	"testing"
)

func Test1(t *testing.T) {
	/*timer.Start(3, 1000, 0, func() {
		time.Sleep(time.Second * 3)
		fmt.Println("xxx")
	})
	c := make(chan int)
	<- c*/

	timer.Start(3, 0, 1000, func() {
		panic("xxx")
	})
	c := make(chan int)
	<-c
}
