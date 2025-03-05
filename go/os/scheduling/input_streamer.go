package main

import (
	"context"
	"sort"
)

type InputStreamer struct {
	tasks             []Task
	TaskInputStream   chan<- Task
	NoMoreInputSignal chan<- bool
}

func (s InputStreamer) InputTask(ctx context.Context) {
	sort.Slice(s.tasks, func(i, j int) bool {
		return s.tasks[i].InputTime < s.tasks[j].InputTime
	})

	for {
		if len(s.tasks) == 0 {
			close(s.NoMoreInputSignal)
			return
		}

		for s.tasks[0].InputTime > int(systemTime.Load()) {
			select {
			case <-ctx.Done():
				panic("Timeout, abort waiting in input streamer")
			default:
				// continue waiting
			}
			// systemTime.Add(1)
		}

		currentInput := []Task{}
		for len(s.tasks) > 0 && s.tasks[0].InputTime <= int(systemTime.Load()) {
			currentInput = append(currentInput, s.tasks[0])
			s.tasks = s.tasks[1:]
		}

		for _, task := range currentInput {
			RecordLog(IO, INPUT_TASK, task)
			s.TaskInputStream <- task
		}
		s.NoMoreInputSignal <- true
	}
}
