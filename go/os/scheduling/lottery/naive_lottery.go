package main

import (
	"fmt"
	"math/rand"
)

const MaxID = 10000
const RangeStart = 0

type naiveLotteryScheduler struct {
	// auto-incrementing one based task ID
	lastId         int
	maxTicketCount int

	taskQueue TaskQueue

	// tracks the number of scheduling per task
	scheduleAudit map[int]int

	logger logger
}

func NewNaiveLotteryScheduler() Scheduler {
	return &naiveLotteryScheduler{
		lastId:         0,
		maxTicketCount: -1,
		logger:         logger{},
		taskQueue:      &sortedTasks{},
		scheduleAudit:  map[int]int{},
	}
}

func (s *naiveLotteryScheduler) ScheduleNextTask() (Schedulable, error) {
	nextTicket := rand.Intn(s.maxTicketCount + 1)

	task, err := s.taskQueue.FindTask(nextTicket)
	if err != nil {
		return nil, err
	}
	s.recordTaskAudit(task)

	return task, nil
}

func (s *naiveLotteryScheduler) AddTask(ticketCount int) Schedulable {
	defer s.increaseMaxTicketCount(ticketCount)

	// assign interval for new task
	var newInterval [2]int
	if len(s.taskQueue.Tasks()) == 0 {
		newInterval = [2]int{RangeStart, RangeStart + ticketCount - 1}
	} else {
		newInterval = [2]int{s.maxTicketCount + 1, s.maxTicketCount + ticketCount}
	}

	// create new task
	s.lastId += 1 % MaxID
	taskId := s.lastId
	newTask := NewSchedulableTask(taskId, ticketCount, newInterval)

	// add to task list
	s.taskQueue.AddTask(newTask)

	s.logger.logTaskAction(newTask, AddTask)

	return newTask
}

func (s *naiveLotteryScheduler) RemoveTask(id int) error {
	task, err := s.taskQueue.RemoveTask(id)
	if err != nil {
		return err
	}
	s.decreaseMaxTicketCount(task.Ticket())

	s.logger.logTaskAction(task, RemoveTask)
	return nil
}

func (s *naiveLotteryScheduler) increaseMaxTicketCount(ticketCount int) {
	s.maxTicketCount += ticketCount
}

func (s *naiveLotteryScheduler) decreaseMaxTicketCount(ticketCount int) {
	s.maxTicketCount -= ticketCount
}

func (s *naiveLotteryScheduler) Log() {
	for _, log := range s.logger.logs {
		fmt.Println(log)
	}
}

func (s *naiveLotteryScheduler) ScheduleAudit() map[int]int {
	return s.scheduleAudit
}

func (s *naiveLotteryScheduler) recordTaskAudit(task Schedulable) {
	if _, exists := s.scheduleAudit[task.ID()]; !exists {
		s.scheduleAudit[task.ID()] = 0
	} else {
		s.scheduleAudit[task.ID()] += 1
	}
}
