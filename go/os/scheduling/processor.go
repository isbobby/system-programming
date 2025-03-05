package main

import (
	"sync/atomic"
)

var systemTime atomic.Int32

type Processor struct {
	TaskQueue          <-chan Task
	CurrentTask        *Task
	Scheduler          Scheduler
	SchedulerCompleted <-chan bool
}

func New() Processor {
	taskChan := make(chan Task)
	schedulerCompletionChan := make(chan bool)

	s := Scheduler{
		TaskDst:            taskChan,
		SchedulerCompleted: schedulerCompletionChan,
	}

	p := Processor{
		TaskQueue:          taskChan,
		SchedulerCompleted: schedulerCompletionChan,
		Scheduler:          s,
	}

	return p
}
