package main

import (
	"fmt"
)

var execStats ExecStats

type ExecStats struct {
	ExecLogs         []ExecLog
	TaskIdToExecLogs map[int][]ExecLog
}

type ExecLog struct {
	Task       Task
	Actor      string
	SystemTime int
	Action     string
}

const (
	START_TASK     = "starting task"
	COMPLETE_TASK  = "complete task"
	SELECT_TASK    = "selected task"
	SCHEDULED_TASK = "scheduled task"
	INPUT_TASK     = "input task"

	PROC = "Processor"
	SCHE = "Scheduler"
	IO   = "IO Stream"
)

func (l *ExecLog) String() string {
	whiteSpace := []byte{}
	if l.SystemTime < 10 {
		whiteSpace = append(whiteSpace, []byte{' ', ' '}...)
	} else if l.SystemTime < 100 {
		whiteSpace = append(whiteSpace, []byte{' '}...)
	}
	return fmt.Sprintf("[time:%v%ds][%v] %v, taskID:%d", string(whiteSpace), l.SystemTime, l.Actor, l.Action, l.Task.Id)
}

func RecordLog(actor string, action string, task Task) {
	if _, exists := execStats.TaskIdToExecLogs[task.Id]; !exists {
		execStats.TaskIdToExecLogs[task.Id] = []ExecLog{}
	}

	newLog := ExecLog{
		Actor:      actor,
		SystemTime: int(systemTime.Load()),
		Action:     action,
		Task:       task,
	}

	execStats.TaskIdToExecLogs[task.Id] = append(execStats.TaskIdToExecLogs[task.Id], newLog)
	execStats.ExecLogs = append(execStats.ExecLogs, newLog)
}

func ShowLog(showExecLog bool) {
	if showExecLog {
		fmt.Println("Exec Log: -------------------------------------------")
		for _, log := range execStats.ExecLogs {
			fmt.Println(log.String())
		}
	}

	fmt.Println("Exec Stats: -----------------------------------------")

	completedCount := 0
	totalTurnAroundTime := 0
	totalResponseTime := 0
	for _, logs := range execStats.TaskIdToExecLogs {
		queueTime := 0
		startTime := 0
		completeTime := 0

		for _, log := range logs {
			queueTime = log.Task.InputTime
			if log.Action == COMPLETE_TASK {
				completeTime = log.SystemTime
			}
			if log.Action == START_TASK {
				startTime = log.SystemTime
			}
			if log.Action == SCHEDULED_TASK {
				startTime = log.SystemTime
			}
		}

		turnAround := completeTime - queueTime
		totalTurnAroundTime += turnAround

		response := startTime - queueTime
		totalResponseTime += response
		// assume all tasks run to completion
		completedCount += 1
	}
	fmt.Printf("[Turn Around] %v / %v = %.2fs\n", totalTurnAroundTime, completedCount, float64(totalTurnAroundTime)/float64(completedCount))
	// fmt.Printf("[Response] %v / %v = %.2f\n", totalResponseTime, completedCount, float64(totalResponseTime)/float64(completedCount))
}
