package main

import "errors"

var (
	ErrIntervalNotExist = errors.New("error due to accessing interval that does not exist")
)

type Scheduler interface {
	ScheduleNextTask() Schedulable

	AddTask(ticketCount int) Schedulable

	RemoveTask(id int) error

	Log()

	ScheduleAudit() map[int]int
}
