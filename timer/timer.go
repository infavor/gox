package timer

import (
	"github.com/hetianyi/gox"
	"github.com/hetianyi/gox/logger"
	"time"
)

// timer defines a timer.
type Timer struct {
	close bool
}

// Destroy destroys the timer.
func (t *Timer) Destroy() {
	t.close = true
}

// Start starts a timer with parameters 'initialDelay', 'fixedDelay', 'fixedRate' and timer work,
// It returns timer struct for controlling.
func Start(initialDelay time.Duration, fixedDelay time.Duration, fixedRate time.Duration, work func(t *Timer)) *Timer {
	logger.Debug("create timer")
	t := &Timer{
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
func (t *Timer) tick(initialDelay time.Duration, fixedDelay time.Duration, fixedRate time.Duration, work func(t *Timer)) {
	logger.Debug("start timer")
	defer func() {
		logger.Debug("stop timer")
	}()
	time.Sleep(initialDelay)
	if t.close {
		return
	}
	if fixedDelay <= 0 {
		tim := time.NewTicker(fixedRate)
		for {
			if t.close {
				tim.Stop()
				break
			}
			gox.Try(func() {
				work(t)
			}, func(i interface{}) {
				logger.Error("error execute timer job:", i)
			})
			<-tim.C
		}
	} else {
		for {
			if t.close {
				break
			}
			gox.Try(func() {
				work(t)
			}, func(i interface{}) {
				logger.Error("error execute timer job:", i)
			})
			time.Sleep(fixedDelay)
		}
	}
}
