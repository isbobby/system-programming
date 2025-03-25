package main

import (
	"context"
	"time"
)

// generate input data for MLFQ
func inputJobs() []*Job {
	initialTime := 0
	return []*Job{
		NewJob(1, initialTime, []JobInput{CPUInstruction, CPUInstruction, CPUInstruction, CPUInstruction}),
		NewJob(2, initialTime, []JobInput{IOInstruction, CPUInstruction, IOInstruction, CPUInstruction}),
		NewJob(3, initialTime, []JobInput{CPUInstruction, CPUInstruction}),
	}
}

func main() {

	clock := NewClock()

	maxPriority := 3
	mlfqResetInterval := 10
	mlfqQueueSize := 100

	sToPChan := make(chan *Job)  // for scheduler to schedule job onto processor, hence S to P
	ioToSChan := make(chan *Job) // for IO to enqueue job for scheduling, hence io To S

	// TODO consider the case when IO buffer is full
	pToIOChan := make(chan *Job, 10) // for processor to swap job for IO, hence p to IO. We have 1 IO device, need a buffer for queue
	pToSChan := make(chan *Job, 1)   // for processor to expire a job back to scheduler, hence p to S

	mlfq := NewMLFQ(
		maxPriority,
		[]int{5, 4, 3, 2},
		mlfqResetInterval,
		mlfqQueueSize,
		clock,
		sToPChan,
		ioToSChan,
		pToSChan,
	)

	initialJobs := inputJobs()
	io := NewIOStream(initialJobs, ioToSChan, pToIOChan)

	ctx := context.Background()
	timeoutCtx, cancel := context.WithTimeout(ctx, 2*time.Second)

	processor := NewProcessor(clock, sToPChan, pToSChan, pToIOChan)

	go io.Run(timeoutCtx, cancel)

	go mlfq.Run(timeoutCtx)

	defer cancel()

	processor.Run(timeoutCtx)
}
