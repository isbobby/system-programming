package main

import "fmt"

type logger struct {
	logs []string
}

type logAction string

var (
	AddTask      logAction = "Add Task"
	ScheduleTask logAction = "Schedule Task"
	UpdateTask   logAction = "Updated Task"
	RemoveTask   logAction = "Removed Task"
)

func (l *logger) logTaskAction(task Schedulable, action logAction, kvs ...string) {
	var log string
	log = fmt.Sprintf("Task (id:%v) - %v", task.ID(), action)
	if len(kvs) > 0 && len(kvs)%2 == 0 {
		kvBytes := []byte(" [")
		for i := 0; i < len(kvs); i += 2 {
			key := kvs[i]
			val := kvs[i+1]
			kvBytes = append(kvBytes, []byte(fmt.Sprintf("%v:%v,", key, val))...)
		}
		kvBytes[len(kvBytes)-1] = ']'
		log = string(append([]byte(log), kvBytes...))
	}
	l.logs = append(l.logs, log)
}
