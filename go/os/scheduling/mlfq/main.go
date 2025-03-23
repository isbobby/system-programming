package main

import (
	"context"
	"time"
)

// generate input data for MLFQ
func inputJobs(maxPriority int) []Job {
	return []Job{
		NewJob(maxPriority, []JobInput{CPUInstruction}),
	}
}

func main() {

	clock := NewClock()

	maxPriority := 3
	mlfqResetInterval := 10
	mlfqQueueSize := 100

	sToPChan := make(chan *Job, 1)  // for scheduler to schedule job onto processor, hence S to P
	ioToSChan := make(chan *Job, 1) // for IO to enqueue job for scheduling, hence io To S
	pToIOChan := make(chan *Job, 1) // for processor to swap job for IO, hence p to IO
	pToSChan := make(chan *Job, 1)  // for processor to expire a job back to scheduler, hence p to S

	mlfq := NewMLFQ(
		maxPriority,
		[]int{4, 3, 2, 1},
		mlfqResetInterval,
		mlfqQueueSize,
		clock,
		sToPChan,
		ioToSChan,
		pToSChan,
	)

	initialJobs := inputJobs(maxPriority)
	io := NewIOStream(initialJobs, ioToSChan, pToIOChan)

	ctx := context.Background()
	timeoutCtx, cancel := context.WithTimeout(ctx, 1*time.Second)

	processor := NewProcessor(clock, sToPChan, pToSChan, pToIOChan)

	go io.ScheduleInput(timeoutCtx, cancel)

	go mlfq.ScheduleJob()

	defer cancel()

	processor.Run(timeoutCtx)
}
