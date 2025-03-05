package main

import (
	"context"
	"sort"
)

var FifoScheduling SchedulingStrategy = func(tasks *[]Task) Task {
	if len(*tasks) == 0 {
		panic("cannot schedule on empty tasks")
	}

	old := *tasks
	sort.Slice(old, func(a, b int) bool {
		return old[a].StartTime < old[b].StartTime
	})

	task := old[0]
	for task.StartTime > int(systemTime.Load()) {
		systemTime.Add(1)
	}
	*tasks = old[1:]

	return task
}

func (p *Processor) RunFifo(ctx context.Context, inputTasks []Task) {
	go p.Scheduler.ScheduleTask(ctx, FifoScheduling, inputTasks)

	for {
		if p.CurrentTask != nil {
			RecordLog(PROC, START_TASK, *p.CurrentTask)
			for p.CurrentTask.Duration > 0 {
				p.CurrentTask.Duration -= 1
				systemTime.Add(1)
			}
			RecordLog(PROC, COMPLETE_TASK, *p.CurrentTask)
			p.CurrentTask = nil
		}

		select {
		case nextTask := <-p.TaskQueue:
			p.CurrentTask = &nextTask
		case <-p.SchedulerCompleted:
		}

		if p.CurrentTask == nil {
			break
		}
	}
	ctx.Done()
}
