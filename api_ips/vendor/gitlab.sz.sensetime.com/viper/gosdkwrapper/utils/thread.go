package utils

import (
	"runtime"

	log "github.com/sirupsen/logrus"
)

type LockedOSThread struct {
	inited   chan bool
	kill     chan bool
	exited   chan bool
	requests chan func()
}

func NewLockedOSThread(threadInitFunc func() error) (*LockedOSThread, error) {
	out := &LockedOSThread{
		inited:   make(chan bool),
		kill:     make(chan bool),
		exited:   make(chan bool),
		requests: make(chan func()),
	}
	go out.background(threadInitFunc)
	log.Info("waiting for thread initialization to finish")
	<-out.inited
	return out, nil
}

func (e *LockedOSThread) background(threadInitFunc func() error) {
	runtime.LockOSThread()
	log.Info("thread initialize started")
	err := threadInitFunc()
	e.inited <- true
	if err != nil {
		panic(err)
	}
	log.Info("thread initialize finished")
	for req := range e.requests {
		req()
	}
	e.exited <- true
	log.Info("locked os thread exit")
}

func (e *LockedOSThread) Execute(f func() interface{}) interface{} {
	done := make(chan interface{}, 1)
	e.requests <- func() {
		ret := f()
		done <- ret
	}
	return <-done
}

func (e *LockedOSThread) Release() {
	close(e.requests)
	<-e.exited
}
