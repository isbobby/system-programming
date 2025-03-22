package main

import (
	"context"
	"sort"
)

var FifoScheduling SchedulingStrategy = func(readyQueue *[]Task) Task {
	if len(*readyQueue) == 0 {
		panic("cannot schedule on empty tasks")
	}

	old := *readyQueue
	sort.Slice(old, func(a, b int) bool {
		return old[a].InputTime < old[b].InputTime
	})

	task := old[0]
	for task.InputTime > int(systemTime.Load()) {
		systemTime.Add(1)
	}
	*readyQueue = old[1:]

	return task
}

func (p *Processor) RunFifo(ctx context.Context, inputTasks []Task) {
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
		case nextTask := <-p.TaskSrcStream:
			p.CurrentTask = &nextTask
		case <-p.NoMoreReadyTaskSignal:
		}

		if p.CurrentTask == nil {
			break
		}
	}
	ctx.Done()
}
