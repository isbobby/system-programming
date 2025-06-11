package main

import (
	"errors"
	"testing"
)

func TestNaiveLotteryAddTasks(t *testing.T) {
	taskTicketCount := []int{1, 3, 5, 4}

	scheduler := naiveLotteryScheduler{
		maxTicketCount: -1,
		taskQueue:      &sortedTasks{},
	}

	for _, count := range taskTicketCount {
		scheduler.AddTask(count)
	}

	// expect sorted interval
	expectedIntervals := [][2]int{{0, 0}, {1, 3}, {4, 8}, {9, 12}}

	for i, intervalToTask := range scheduler.taskQueue.Tasks() {
		if intervalToTask.Interval() != expectedIntervals[i] {
			t.Error("interval result differs from expected", "expects", expectedIntervals[i], "got", intervalToTask.Interval())
			return
		}
	}

	if scheduler.maxTicketCount != 12 {
		t.Error("max ticket count does not match expectation", "expects", 12, "got", scheduler.maxTicketCount)
	}
}

func TestNaiveLotteryRemoveTaskDoesNotExist(t *testing.T) {
	scheduler := NewNaiveLotteryScheduler()
	err := scheduler.RemoveTask(1)
	if err == nil {
		t.Error("expects scheduler to return err")
	}
	target := ErrIntervalNotExist
	if !errors.Is(err, target) {
		t.Error("expects scheduler to return ErrIntervalNotExist")
	}
}

func TestNaiveLotteryRemoveTasks(t *testing.T) {
	taskTicketCount := []int{1, 3, 5, 4}

	scheduler := naiveLotteryScheduler{
		maxTicketCount: -1,
		taskQueue:      &sortedTasks{},
	}

	for _, count := range taskTicketCount {
		scheduler.AddTask(count)
	}
	err := scheduler.RemoveTask(2)
	if err != nil {
		t.Error("expects no error from removal of valid task")
		return
	}

	// expect compacted sorted interval
	// original {{0, 0}, {1, 3}, {4, 8}, {9, 12}}
	expectedIntervals := [][2]int{{0, 0}, {1, 5}, {6, 9}}
	if len(expectedIntervals) != len(scheduler.taskQueue.Tasks()) {
		t.Error("task queue size differs from expected", "expects", len(expectedIntervals), "got", scheduler.taskQueue.Tasks())
		return
	}

	for i, task := range scheduler.taskQueue.Tasks() {
		if task.Interval() != expectedIntervals[i] {
			t.Error("interval result differs from expected", "expects", expectedIntervals[i], "got", task.Interval())
			return
		}
	}

	if scheduler.maxTicketCount != 9 {
		t.Error("max ticket count does not match expectation", "expectes", 9, "got", scheduler.maxTicketCount)
	}
}

func TestNaiveLotteryRemoveLastTask(t *testing.T) {
	taskTicketCount := []int{1, 3, 5, 4}

	scheduler := naiveLotteryScheduler{
		maxTicketCount: -1,
		taskQueue:      &sortedTasks{},
	}

	for _, count := range taskTicketCount {
		scheduler.AddTask(count)
	}
	err := scheduler.RemoveTask(4)
	if err != nil {
		t.Error("expects no error from removal of valid task")
		return
	}

	// expect compacted sorted interval
	expectedIntervals := [][2]int{{0, 0}, {1, 3}, {4, 8}}
	if len(expectedIntervals) != len(scheduler.taskQueue.Tasks()) {
		t.Error("interval length differs from expected", "expects", len(expectedIntervals), "got", len(scheduler.taskQueue.Tasks()))
		return
	}

	for i, task := range scheduler.taskQueue.Tasks() {
		if task.Interval() != expectedIntervals[i] {
			t.Error("interval result differs from expected", "expects", expectedIntervals[i], "got", task.Interval())
			return
		}
	}

	if scheduler.maxTicketCount != 8 {
		t.Error("max ticket count does not match expectation", "expectes", 8, "got", scheduler.maxTicketCount)
	}
}

func TestNaiveLotteryRemoveFirstTask(t *testing.T) {
	taskTicketCount := []int{1, 3, 5, 4}

	scheduler := naiveLotteryScheduler{
		maxTicketCount: -1,
		taskQueue:      &sortedTasks{},
	}

	for _, count := range taskTicketCount {
		scheduler.AddTask(count)
	}
	err := scheduler.RemoveTask(1)
	if err != nil {
		t.Error("expects no error from removal of valid task")
		return
	}

	// expect compacted sorted interval
	expectedIntervals := [][2]int{{0, 2}, {3, 7}, {8, 11}}
	if len(expectedIntervals) != len(scheduler.taskQueue.Tasks()) {
		t.Error("interval length differs from expected", "expects", len(expectedIntervals), "got", len(scheduler.taskQueue.Tasks()))
		return
	}

	for i, intervalToTask := range scheduler.taskQueue.Tasks() {
		if intervalToTask.Interval() != expectedIntervals[i] {
			t.Error("interval result differs from expected", "expects", expectedIntervals[i], "got", intervalToTask.Interval())
			return
		}
	}

	if scheduler.maxTicketCount != 11 {
		t.Error("max ticket count does not match expectation", "expectes", 11, "got", scheduler.maxTicketCount)
	}
}

func TestSchedulingApproachesFairShare(t *testing.T) {
	taskTicketCount := []int{1, 2, 2, 3}
	tolerance := 0.02

	scheduler := NewNaiveLotteryScheduler()

	for _, count := range taskTicketCount {
		scheduler.AddTask(count)
	}
	scheduler.RemoveTask(2)
	scheduler.AddTask(4)

	expectedShare := map[int]float64{
		1: 0.1,
		3: 0.2,
		4: 0.3,
		5: 0.4,
	}

	totalSchedules := 10000
	for i := 0; i < totalSchedules; i++ {
		scheduler.ScheduleNextTask()
	}

	statsByTaskId := scheduler.ScheduleAudit()

	for task, expectedShare := range expectedShare {
		actualShare := float64(statsByTaskId[task]) / float64(totalSchedules)

		if err := acceptableProbability(expectedShare, actualShare, tolerance); err != nil {
			t.Error("task share not within tolerable range", "task expected", expectedShare, "task actual", actualShare)
			t.Log(scheduler.ScheduleAudit())
			scheduler.Log()
			return
		}
	}
}
