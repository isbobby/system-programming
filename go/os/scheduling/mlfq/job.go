package main

import "sync/atomic"

type Job struct {
	ScheduledTime int
	Priority      int
	Inputs        []JobInput
	TimeAllotment atomic.Int32
}

type JobInput struct {
	Cycle int
	Type  string
}

var (
	IOInstruction JobInput = JobInput{
		Cycle: 5,
		Type:  "IO",
	}
	CPUInstruction JobInput = JobInput{
		Cycle: 1,
		Type:  "CPU",
	}
)

func NewJob(maxPriority int, instructions []JobInput) Job {
	return Job{
		Priority: maxPriority,
		Inputs:   instructions,
	}
}

func (j *Job) DecreasePriority() {
	if j.Priority == 0 {
		return
	}

	j.Priority -= 1
}
