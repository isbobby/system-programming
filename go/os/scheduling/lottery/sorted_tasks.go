package main

type sortedTasks []Schedulable

func (s *sortedTasks) Tasks() []Schedulable {
	return *s
}

func (s *sortedTasks) RemoveTask(id int) (Schedulable, error) {
	sortedTaskSlice := *s
	index := -1
	for i := range *s {
		if sortedTaskSlice[i].ID() == id {
			index = i
			break
		}
	}
	if index < 0 {
		return nil, ErrIntervalNotExist
	}

	// remove last task
	// refactor this, because if err is introduced below, we may get unwanted update since defer is not conditional
	removedTask := sortedTaskSlice[index]
	if index == len(sortedTaskSlice)-1 {
		*s = sortedTaskSlice[:len(sortedTaskSlice)-1]
		return removedTask, nil
	}

	// remove task in other position
	previousIndex := index - 1
	var shiftedIntervalStart int
	if previousIndex < 0 {
		// if removed task is the head
		shiftedIntervalStart = 0
	} else {
		shiftedIntervalStart = sortedTaskSlice[previousIndex].Interval()[1] + 1
	}

	for i := index + 1; i < len(sortedTaskSlice); i++ {
		oldTaskInterval := sortedTaskSlice[i].Interval()
		taskTicketCount := oldTaskInterval[1] - oldTaskInterval[0]
		newInterval := [2]int{shiftedIntervalStart, shiftedIntervalStart + taskTicketCount}
		sortedTaskSlice[i].SetInterval(newInterval)
		shiftedIntervalStart = newInterval[1] + 1

		// s.logger.logTaskAction(
		// 	s.sortedTaskList[i], UpdateTask,
		// 	"old interval", fmt.Sprintf("[%v,%v]", oldTaskInterval[0], oldTaskInterval[1]),
		// 	"new interval", fmt.Sprintf("[%v,%v]", newInterval[0], newInterval[1]),
		// )
	}

	newSortedTaskList := sortedTaskSlice[:index]
	newSortedTaskList = append(newSortedTaskList, sortedTaskSlice[index+1:]...)
	sortedTaskSlice = newSortedTaskList
	*s = sortedTaskSlice
	return removedTask, nil
}

func (s *sortedTasks) AddTask(task Schedulable) error {
	*s = append(*s, task)
	return nil
}

func (s *sortedTasks) FindTask(ticket int) (Schedulable, error) {
	tasks := *s

	if len(tasks) == 0 {
		return nil, ErrIntervalNotExist
	}

	start, end := 0, len(tasks)-1

	for start < end {
		mid := (start + end) / 2

		if tasks[mid].MatchInterval(ticket) {
			return tasks[mid], nil
		}

		if ticket < tasks[mid].Interval()[0] {
			end = mid - 1
		} else {
			start = mid + 1
		}
	}

	return tasks[end], nil
}
