package main

import (
	"context"
	"fmt"
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
	if task.StartTime > int(systemTime.Load()) {
		systemTime.Add(1)
	}
	*tasks = old[1:]

	return task
}

func (p Processor) RunFifo(ctx context.Context, inputTasks []Task) {
	go p.Scheduler.ScheduleTask(ctx, FifoScheduling, inputTasks)

	for !p.SchedulerReady.Load() {
		// wait for scheduler to be ready
	}

	for {
		if p.CurrentTask != nil {
			fmt.Println("Executing", p.CurrentTask.Id, "t", systemTime.Load())
			for p.CurrentTask.Duration > 0 {
				p.CurrentTask.Duration -= 1
				systemTime.Add(1)
			}
			fmt.Println("Completed", p.CurrentTask.Id, "t", systemTime.Load())
			p.CurrentTask = nil
		}

		select {
		case nextTask := <-p.TaskQueue:
			p.CurrentTask = &nextTask
			fmt.Println("Received", p.CurrentTask.Id, "t", systemTime.Load())
		default:
			// no more task
			fmt.Println("No more task to execute")
		}

		if p.CurrentTask == nil {
			break
		}
	}
	ctx.Done()
}
