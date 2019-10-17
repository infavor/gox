package disk

import (
	"container/list"
	"github.com/hetianyi/gox"
	"github.com/hetianyi/gox/logger"
	"sync"
)

type Entry struct {
	Source  interface{}
	ErrChan chan error
}

type SwitchBuffer struct {

	// left bucket buffer(default bucket).
	left *list.List

	// right bucket buffer.
	right *list.List

	// current buffer list.
	currentBuffer *list.List

	handler func(items *list.List) error

	lock *sync.Mutex

	// mark current status is left or right,
	// when it is true, current buffer list is Left.
	flag bool

	busy bool

	scheduleChan chan bool
	scheduleLock *sync.Mutex

	close bool
}

// NewSwitchBuffer create a new SwitchBuffer.
func NewSwitchBuffer(handler func(items *list.List) error) *SwitchBuffer {
	ret := &SwitchBuffer{
		left:         list.New(),
		right:        list.New(),
		lock:         new(sync.Mutex),
		flag:         true,
		handler:      handler,
		scheduleChan: make(chan bool),
		scheduleLock: new(sync.Mutex),
		close:        false,
	}
	ret.currentBuffer = ret.left
	return ret
}

// Switch switch buffer bucket and clear the buffer.
func (s *SwitchBuffer) switchFlag() {
	s.flag = !s.flag
	s.currentBuffer = gox.TValue(s.flag, s.left, s.right).(*list.List)
	for ele := s.getWorkBucket().Front(); ele != nil; ele = ele.Next() {
		s.currentBuffer.Remove(ele)
	}
}

// Cache gets current writable buffer.
func (s *SwitchBuffer) Push(item interface{}) chan error {
	s.lock.Lock()
	defer s.lock.Unlock()
	errC := make(chan error)
	s.currentBuffer.PushBack(&Entry{
		Source:  item,
		ErrChan: errC,
	})
	if s.currentBuffer.Len() == 0 {
		s.scheduleChan <- true
	}
	return errC
}

// Cache gets current writable buffer.
func (s *SwitchBuffer) Destroy() {
	s.close = true
	s.scheduleChan <- true
}

// must run in a personal goroutine.
func (s *SwitchBuffer) Schedule() {
	for {
		logger.Info("Schedule is waiting..")
		for s.currentBuffer.Len() == 0 {
			<-s.scheduleChan
		}
		if s.busy {
			continue
		}
		s.scheduleLock.Lock()
		if s.close {
			break
		}
		logger.Info("Schedule running..")
		// if current buffer has jobs, then handle them.
		if s.currentBuffer.Len() > 0 {
			// switch flag and handle jobs.
			s.switchFlag()
			s.work()
		}
		s.scheduleLock.Unlock()
	}
}

func (s *SwitchBuffer) getWorkBucket() *list.List {
	return gox.TValue(s.flag, s.right, s.left).(*list.List)
}

func (s *SwitchBuffer) work() {
	defer func() {
		go func() {
			logger.Info("chan 2--------------------")
			s.scheduleChan <- true
			logger.Info("chan 2--------------------passed")
		}()
	}()
	bc := s.getWorkBucket()
	err := s.handler(bc)
	gox.WalkList(bc, func(item interface{}) bool {
		entry := item.(*Entry)
		logger.Info("chan 3--------------------")
		entry.ErrChan <- err
		logger.Info("chan 3--------------------passed")
		return false
	})
	logger.Info("job end--------------------")
}
