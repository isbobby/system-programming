package main

import (
	"context"
)

func main() {
	cases := [][]Task{caseOne, caseTwo}
	for i := range cases {
		// reset machine state
		execStats = ExecStats{
			ExecLogs:         []ExecLog{},
			TaskIdToExecLogs: map[int][]ExecLog{},
		}
		newProcessor := New()
		systemTime.Store(0)
		ctx := context.Background()

		newProcessor.RunFifo(ctx, cases[i])
		ShowLog(true)
	}
}
