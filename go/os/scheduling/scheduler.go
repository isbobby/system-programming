package main

import (
	"context"
)

type SchedulingStrategy func(tasks *[]Task) Task

type Scheduler struct {
	Queue              []Task
	TaskDst            chan<- Task
	SwtichQueue        <-chan Task
	SchedulerCompleted chan bool
}

func (s Scheduler) ScheduleTask(ctx context.Context, strategy SchedulingStrategy, tasks []Task) {
	s.Queue = tasks
	for {
		if len(s.Queue) != 0 {
			// select a task using scheduling strategy
			nextTask := strategy(&s.Queue)
			RecordLog(SCHE, SELECT_TASK, nextTask)
			s.TaskDst <- nextTask
			RecordLog(SCHE, SCHEDULED_TASK, nextTask)
		}

		if len(s.Queue) == 0 {
			// no more task
			close(s.SchedulerCompleted)
			return
		}

		select {
		case swappedTask := <-s.SwtichQueue:
			s.Queue = append(s.Queue, swappedTask)
		default:
		}
	}
}
