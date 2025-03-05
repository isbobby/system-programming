package main

import (
	"sync/atomic"
)

var systemTime atomic.Int32

type Processor struct {
	TaskSrcStream         <-chan Task
	NoMoreReadyTaskSignal <-chan bool

	CurrentTask *Task
}

func New(TaskSrcStream <-chan Task, schedulerCompletionChan <-chan bool) Processor {
	p := Processor{
		TaskSrcStream:         TaskSrcStream,
		NoMoreReadyTaskSignal: schedulerCompletionChan,
	}
	return p
}
