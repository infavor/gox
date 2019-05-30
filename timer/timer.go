package timer

import (
	"github.com/hetianyi/gox"
	"github.com/hetianyi/gox/logger"
	"github.com/sirupsen/logrus"
	"time"
)

func init() {
	logger.Init(nil)
}

// timer defines a timer.
type timer struct {
	close bool
}

// Destroy destroys the timer.
func (timer *timer) Destroy() {
	timer.close = true
}

// Start starts a timer with parameters 'initialDelay', 'fixedDelay', 'fixedRate' and timer work,
// It returns timer struct for controlling.
func Start(initialDelay time.Duration, fixedDelay time.Duration, fixedRate time.Duration, work func()) *timer {
	t := &timer{
		close: false,
	}
	go t.tick(initialDelay, fixedDelay, fixedRate, work)
	return t
}

// Timer defines a simple way to use is a timer.
// initialDelay: number of milliseconds to delay before the first call.
// fixedDelay: a fixed period in milliseconds between the
// end of the last call and the start of the next.
// fixedRate: a fixed period in milliseconds between calls.
//
// Note that fixedDelay is superior than fixedRate.
func (timer *timer) tick(initialDelay time.Duration, fixedDelay time.Duration, fixedRate time.Duration, work func()) {
	time.Sleep(initialDelay)
	if timer.close {
		return
	}
	if fixedDelay <= 0 {
		t := time.NewTicker(fixedRate)
		for {
			if timer.close {
				break
			}
			gox.Try(func() {
				work()
			}, func(i interface{}) {
				logrus.Error("error execute timer job:", i)
			})
			<-t.C
		}
	} else {
		for {
			if timer.close {
				break
			}
			gox.Try(func() {
				work()
			}, func(i interface{}) {
				logrus.Error("error execute timer job:", i)
			})
			time.Sleep(fixedDelay)
		}
	}
}
