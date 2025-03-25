package main

import "sync/atomic"

type Job struct {
	ID               int
	ScheduledTime    int
	Priority         *int
	InstructionStack []JobInput
	TimeAllotment    atomic.Int32
}

type JobInput struct {
	Cycle int
	Type  string
}

func (i JobInput) IsIO() bool {
	return i.Type == "IO"
}

func (i JobInput) IsCPU() bool {
	return i.Type == "CPU"
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

func NewJob(ID int, InstructionStack []JobInput) *Job {
	return &Job{
		ID:               ID,
		InstructionStack: InstructionStack,
	}
}

func (j *Job) DecreasePriority() {
	if *j.Priority == 0 {
		return
	}

	*j.Priority -= 1
}
