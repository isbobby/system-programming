package main

import (
	"context"
	"fmt"
	"sync/atomic"
)

type SchedulingStrategy func(tasks *[]Task) Task

type Scheduler struct {
	Queue          []Task
	TaskDst        chan<- Task
	SwtichQueue    <-chan Task
	SchedulerReady *atomic.Bool
}

func (s Scheduler) ScheduleTask(ctx context.Context, strategy SchedulingStrategy, tasks []Task) {
	s.Queue = tasks

	for {
		if len(s.Queue) != 0 {
			nextTask := strategy(&s.Queue)
			fmt.Println("selected task", nextTask, "task left", s.Queue)
			s.TaskDst <- nextTask
			s.SchedulerReady.Store(true)
		}

		select {
		case swappedTask := <-s.SwtichQueue:
			s.Queue = append(s.Queue, swappedTask)
		default:
		}
	}
}
