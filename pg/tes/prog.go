package main

import (
	"github.com/hetianyi/gox/pg"
	"time"
)

func main() {
	p := pg.New(100, 50, "False Alarm.mp3", pg.Left)
	for i := 0; i < 100; i++ {
		p.Increase()
		time.Sleep(time.Millisecond * 100)
	}
	time.Sleep(time.Second * 10)
}
