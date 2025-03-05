package main

import (
	"context"
)

type SchedulingStrategy func(tasks *[]Task) Task

type Scheduler struct {
	readyTasks []Task

	TaskDstStream         chan<- Task
	TaskSwtichStream      <-chan Task // TODO
	TaskInputStream       <-chan Task
	NoMoreReadyTaskSignal chan bool
	NoMoreInputSignal     chan bool
}

func NewScheduler() Scheduler {
	return Scheduler{}
}

func (s Scheduler) ScheduleTask(ctx context.Context, strategy SchedulingStrategy) {
	for {
		if len(s.readyTasks) != 0 {
			// select a task using scheduling strategy
			nextTask := strategy(&s.readyTasks)
			RecordLog(SCHE, SELECT_TASK, nextTask)
			s.TaskDstStream <- nextTask
			RecordLog(SCHE, SCHEDULED_TASK, nextTask)
		}

	switchChanLoop:
		for {
			select {
			case swappedTask := <-s.TaskSwtichStream:
				s.readyTasks = append(s.readyTasks, swappedTask)
			default:
				break switchChanLoop
			}
		}

	inputChanLoop:
		for {
			select {
			case inputTask := <-s.TaskInputStream:
				s.readyTasks = append(s.readyTasks, inputTask)
			case <-s.NoMoreInputSignal:
				break inputChanLoop
			}
		}

		if len(s.readyTasks) == 0 {
			close(s.NoMoreReadyTaskSignal)
			return
		}
	}
}
