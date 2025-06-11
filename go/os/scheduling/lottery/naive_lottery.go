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
	sortedTaskList []Schedulable

	// tracks the number of scheduling per task
	scheduleAudit map[int]int

	logger logger
}

func NewNaiveLotteryScheduler() Scheduler {
	return &naiveLotteryScheduler{
		lastId:         0,
		maxTicketCount: 0,
		logger:         logger{},
		sortedTaskList: []Schedulable{},
		scheduleAudit:  map[int]int{},
	}
}

func (s *naiveLotteryScheduler) ScheduleNextTask() Schedulable {
	nextTicket := rand.Intn(s.maxTicketCount + 1)

	nextTaskIndex := s.searchTaskByTicketNumber(nextTicket)

	scheduledTask := s.sortedTaskList[nextTaskIndex]

	s.logger.logTaskAction(scheduledTask, ScheduleTask)

	if _, exists := s.scheduleAudit[scheduledTask.ID()]; !exists {
		s.scheduleAudit[scheduledTask.ID()] = 1
	}
	s.scheduleAudit[scheduledTask.ID()] += 1

	return scheduledTask
}

func (s *naiveLotteryScheduler) AddTask(ticketCount int) Schedulable {
	defer s.updateMaxTicketCount()

	// assign interval for new task
	var newInterval [2]int
	if len(s.sortedTaskList) == 0 {
		newInterval = [2]int{RangeStart, RangeStart + ticketCount - 1}
	} else {
		lastInterval := s.sortedTaskList[len(s.sortedTaskList)-1].Interval()
		newInterval = [2]int{lastInterval[1] + 1, lastInterval[1] + ticketCount}
	}

	// create new task
	s.lastId += 1 % MaxID
	taskId := s.lastId
	newTask := NewSchedulableTask(taskId, ticketCount, newInterval)

	// add to task list
	s.sortedTaskList = append(s.sortedTaskList, newTask)
	s.logger.logTaskAction(newTask, AddTask)

	return newTask
}

func (s *naiveLotteryScheduler) RemoveTask(id int) error {
	index := -1
	for i := range s.sortedTaskList {
		if s.sortedTaskList[i].ID() == id {
			index = i
			break
		}
	}
	if index < 0 {
		return ErrIntervalNotExist
	}

	s.logger.logTaskAction(s.sortedTaskList[index], RemoveTask)

	// remove last task
	// refactor this, because if err is introduced below, we may get unwanted update since defer is not conditional
	defer s.updateMaxTicketCount()
	if index == len(s.sortedTaskList)-1 {
		s.sortedTaskList = s.sortedTaskList[:len(s.sortedTaskList)-1]
		return nil
	}

	// remove task in other position
	previousIndex := index - 1
	var shiftedIntervalStart int
	if previousIndex < 0 {
		// if removed task is the head
		shiftedIntervalStart = 0
	} else {
		shiftedIntervalStart = s.sortedTaskList[previousIndex].Interval()[1] + 1
	}

	for i := index + 1; i < len(s.sortedTaskList); i++ {
		oldTaskInterval := s.sortedTaskList[i].Interval()
		taskTicketCount := oldTaskInterval[1] - oldTaskInterval[0]
		newInterval := [2]int{shiftedIntervalStart, shiftedIntervalStart + taskTicketCount}
		s.sortedTaskList[i].SetInterval(newInterval)
		shiftedIntervalStart = newInterval[1] + 1

		s.logger.logTaskAction(
			s.sortedTaskList[i], UpdateTask,
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
	s.maxTicketCount = s.sortedTaskList[len(s.sortedTaskList)-1].Interval()[1]
}

func (s *naiveLotteryScheduler) searchTaskByTicketNumber(t int) int {
	if len(s.sortedTaskList) == 1 {
		return 0
	}

	start, end := 0, len(s.sortedTaskList)-1

	for start < end {
		mid := (start + end) / 2

		currInterval := s.sortedTaskList[mid].Interval()
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
