package timer

import "time"

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
func Start(initialDelay int64, fixedDelay int64, fixedRate int64, work func()) *timer {
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
func (timer *timer) tick(initialDelay int64, fixedDelay int64, fixedRate int64, work func()) {
	time.Sleep(time.Millisecond * time.Duration(initialDelay))
	if timer.close {
		return
	}
	if fixedDelay <= 0 {
		t := time.NewTicker(time.Millisecond * time.Duration(fixedRate))
		for {
			if timer.close {
				break
			}
			func() {
				defer func() {
					if err := recover(); err != nil {
						// do nothing
					}
				}()
				work()
			}()
			<-t.C
		}
	} else {
		for {
			if timer.close {
				break
			}
			func() {
				defer func() {
					if err := recover(); err != nil {
						// do nothing
					}
				}()
				work()
			}()
			time.Sleep(time.Millisecond * time.Duration(fixedDelay))
		}
	}
}
