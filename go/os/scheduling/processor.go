package main

import (
	"sync/atomic"
)

var systemTime atomic.Int32

type Processor struct {
	TaskQueue      <-chan Task
	SchedulerReady *atomic.Bool
	Scheduler      Scheduler
	CurrentTask    *Task
}

func New() Processor {
	taskChan := make(chan Task, 1)

	schedulerReady := atomic.Bool{}
	schedulerReady.Store(false)

	scheduler := Scheduler{
		TaskDst:        taskChan,
		SchedulerReady: &schedulerReady,
	}
	return Processor{
		Scheduler:      scheduler,
		TaskQueue:      taskChan,
		SchedulerReady: &schedulerReady,
	}
}
