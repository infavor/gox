// Copyright (C) 2019 tisnyo <tisnyo@gmail.com>.
//
// A pool allows a particular type of work to run in it and limits its parallel number and maximum wait queue length.
// license that can be found in the LICENSE file.
package pool

import (
	"container/list"
	"errors"
	"github.com/infavor/gox"
	"github.com/infavor/gox/logger"
	"sync"
)

type WorkPool interface {
	Push(task func()) error
}

// pool is a task pool which can limit the number of concurrent task.
type pool struct {
	coreSize          int
	maxWait           int
	activeTaskSize    int
	listPushLock      *sync.Mutex
	listOperationLock *sync.Mutex
	listFetchLock     *sync.Mutex
	numLock           *sync.Mutex
	cha               chan int
	waitingList       *list.List
}

// New creates a task pool.
func New(coreSize int, maxWait int) WorkPool {
	p := &pool{
		coreSize:          coreSize,
		maxWait:           maxWait,
		activeTaskSize:    0,
		listPushLock:      new(sync.Mutex),
		listFetchLock:     new(sync.Mutex),
		listOperationLock: new(sync.Mutex),
		numLock:           new(sync.Mutex),
		cha:               make(chan int),
		waitingList:       list.New(),
	}
	go p.taskWatcher()
	return p
}

// Push push a new task into waiting list
func (pool *pool) Push(task func()) error {
	pool.listPushLock.Lock()
	defer pool.listPushLock.Unlock()
	if pool.waitingList.Len() == pool.maxWait {
		return errors.New("pool is full")
	}
	pool.listOperation(true, task)
	if pool.waitingList.Len() > 0 && pool.updateActiveTaskSize(0) < pool.coreSize {
		pool.cha <- 1
	}
	return nil
}

// listOperation ensures operating on list.List is the only one place.
func (pool *pool) listOperation(push bool, work func()) func() {
	pool.listOperationLock.Lock()
	defer pool.listOperationLock.Unlock()
	if push {
		pool.waitingList.PushBack(work)
		return nil
	}
	if pool.waitingList.Len() > 0 {
		return pool.waitingList.Remove(pool.waitingList.Front()).(func())
	}
	return nil
}

// Push push a new task into waiting list
func (pool *pool) fetchTask() func() {
	pool.listFetchLock.Lock()
	defer pool.listFetchLock.Unlock()
	for pool.waitingList.Len() == 0 || pool.updateActiveTaskSize(0) >= pool.coreSize {
		<-pool.cha
	}
	return pool.listOperation(false, nil)
}

// updateActiveTaskSize update current pool's active task size.
func (pool *pool) updateActiveTaskSize(increment int) int {
	pool.numLock.Lock()
	defer pool.numLock.Unlock()
	pool.activeTaskSize += increment
	return pool.activeTaskSize
}

// taskWatcher watches task list and executes them.
func (pool *pool) taskWatcher() {
	for {
		task := pool.fetchTask()
		if task != nil {
			pool.updateActiveTaskSize(1)
			go pool.execute(task)
		}
	}
}

// execute execute a task and callback function.
func (pool *pool) execute(task func()) {
	defer func() {
		pool.updateActiveTaskSize(-1)
		pool.cha <- 1
	}()
	gox.Try(func() {
		task()
	}, func(i interface{}) {
		logger.Error("error execute work:", i)
	})
}
