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
	sortedTaskList []intervalToTask

	// tracks the number of scheduling per task
	scheduleAudit map[int]int

	logger logger
}

type intervalToTask struct {
	interval [2]int
	task     Schedulable
}

func NewNaiveLotteryScheduler() Scheduler {
	return &naiveLotteryScheduler{
		lastId:         0,
		maxTicketCount: 0,
		logger:         logger{},
		sortedTaskList: []intervalToTask{},
		scheduleAudit:  map[int]int{},
	}
}

func (s *naiveLotteryScheduler) ScheduleNextTask() Schedulable {
	nextTicket := rand.Intn(s.maxTicketCount + 1)

	nextTaskIndex := s.searchTaskByTicketNumber(nextTicket)

	scheduledTask := s.sortedTaskList[nextTaskIndex]

	s.logger.logTaskAction(scheduledTask.task, ScheduleTask)

	if _, exists := s.scheduleAudit[scheduledTask.task.ID()]; !exists {
		s.scheduleAudit[scheduledTask.task.ID()] = 1
	}
	s.scheduleAudit[scheduledTask.task.ID()] += 1

	return scheduledTask.task
}

func (s *naiveLotteryScheduler) AddTask(ticketCount int) Schedulable {
	defer s.updateMaxTicketCount()

	// create new task
	s.lastId += 1 % MaxID
	taskId := s.lastId

	newTask := NewTask(taskId, ticketCount)
	defer s.logger.logTaskAction(newTask, AddTask)

	// update ticket range mapping
	var newInterval [2]int
	if len(s.sortedTaskList) == 0 {
		newInterval = [2]int{RangeStart, RangeStart + ticketCount - 1}
	} else {
		lastInterval := s.sortedTaskList[len(s.sortedTaskList)-1].interval
		newInterval = [2]int{lastInterval[1] + 1, lastInterval[1] + ticketCount}
	}

	newIntervalToTask := intervalToTask{
		interval: newInterval,
		task:     newTask,
	}
	s.sortedTaskList = append(s.sortedTaskList, newIntervalToTask)

	return newTask
}

func (s *naiveLotteryScheduler) RemoveTask(id int) error {
	defer s.updateMaxTicketCount()

	index := -1
	for i := range s.sortedTaskList {
		if s.sortedTaskList[i].task.ID() == id {
			index = i
			break
		}
	}
	if index < 0 {
		return ErrIntervalNotExist
	}

	s.logger.logTaskAction(s.sortedTaskList[index].task, RemoveTask)

	// remove last task
	if index == len(s.sortedTaskList)-1 {
		s.sortedTaskList = s.sortedTaskList[:len(s.sortedTaskList)-1]
		return nil
	}

	previousIndex := index - 1
	var shiftedIntervalStart int
	if previousIndex < 0 {
		// if removed task is the head
		shiftedIntervalStart = 0
	} else {
		shiftedIntervalStart = s.sortedTaskList[previousIndex].interval[1] + 1
	}

	for i := index + 1; i < len(s.sortedTaskList); i++ {
		oldTaskInterval := s.sortedTaskList[i].interval
		taskTicketCount := oldTaskInterval[1] - oldTaskInterval[0]
		newInterval := [2]int{shiftedIntervalStart, shiftedIntervalStart + taskTicketCount}
		s.sortedTaskList[i].interval = newInterval
		shiftedIntervalStart = newInterval[1] + 1

		s.logger.logTaskAction(
			s.sortedTaskList[i].task, UpdateTask,
			"old interval", fmt.Sprintf("[%v,%v]", oldTaskInterval[0], oldTaskInterval[1]),
			"new interval", fmt.Sprintf("[%v,%v]", newInterval[0], newInterval[1]),
		)
	}

	newSortedTaskList := s.sortedTaskList[:index]
	newSortedTaskList = append(newSortedTaskList, s.sortedTaskList[index+1:]...)
	s.sortedTaskList = newSortedTaskList

	return nil
}

func (s *naiveLotteryScheduler) updateMaxTicketCount() {
	s.maxTicketCount = s.sortedTaskList[len(s.sortedTaskList)-1].interval[1]
}

func (s *naiveLotteryScheduler) searchTaskByTicketNumber(t int) int {
	if len(s.sortedTaskList) == 1 {
		return 0
	}

	start, end := 0, len(s.sortedTaskList)-1

	for start < end {
		mid := (start + end) / 2

		currInterval := s.sortedTaskList[mid].interval
		if currInterval[0] <= t && t <= currInterval[1] {
			return mid
		}

		if t < currInterval[0] {
			end = mid - 1
		} else {
			start = mid + 1
		}
	}

	return end
}

func (s *naiveLotteryScheduler) Log() {
	for _, log := range s.logger.logs {
		fmt.Println(log)
	}
}

func (s naiveLotteryScheduler) ScheduleAudit() map[int]int {
	return s.scheduleAudit
}
