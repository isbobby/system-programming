package main

type Job struct {
	Priority int
	Inputs   []JobInput
}

type JobInput string

var (
	IOInstruction  JobInput = "IO"
	CPUInstruction JobInput = "CPU"
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
