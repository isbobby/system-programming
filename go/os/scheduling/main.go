package main

import (
	"context"
	"time"
)

func main() {
	cases := [][]Task{caseOne, caseTwo}
	for i := range cases {
		// reset machine state
		execStats = ExecStats{
			ExecLogs:         []ExecLog{},
			TaskIdToExecLogs: map[int][]ExecLog{},
		}

		taskInputStream := make(chan Task)
		noMoreInputSignal := make(chan bool)
		input := InputStreamer{
			tasks:             cases[i],
			TaskInputStream:   taskInputStream,
			NoMoreInputSignal: noMoreInputSignal,
		}

		noMoreReadyTaskSignal := make(chan bool)
		taskScheduleStream := make(chan Task)
		taskSwitchStream := make(chan Task)
		s := Scheduler{
			TaskDstStream:         taskScheduleStream,
			TaskInputStream:       taskInputStream,
			TaskSwtichStream:      taskSwitchStream,
			NoMoreReadyTaskSignal: noMoreReadyTaskSignal,
			NoMoreInputSignal:     noMoreInputSignal,
		}

		p := Processor{
			TaskSrcStream:         taskScheduleStream,
			NoMoreReadyTaskSignal: noMoreReadyTaskSignal,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		systemTime.Store(0)

		go input.InputTask(ctx)
		go s.ScheduleTask(ctx, FifoScheduling)
		p.RunFifo(ctx, cases[i])
		ShowLog(true)
	}
}
