package main

import (
	"context"
	"time"
)

func RunSystem(MLFQCfg *MLFQConfig, inputs []*Job, timeout time.Duration, verbose bool) []AuditLog {
	ctx := context.Background()
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	clockToIOSignal := make(chan interface{})
	clockToPSignal := make(chan interface{})

	clockSubscriptions := []chan<- interface{}{clockToIOSignal, clockToPSignal}

	clock := NewClock(time.Duration(250*time.Millisecond), clockSubscriptions)

	pToSSignal := make(chan interface{})
	sToIOSignal := make(chan interface{})
	ioCompletedSignal := make(chan interface{})

	ioToSChan := make(chan *Job, 10) // for IO to enqueue job for scheduling, hence io To S, assuming 10 job buffer
	pToSChan := make(chan *Job, 1)   // for processor to expire a job back to scheduler, hence p to S
	sToPChan := make(chan *Job)      // for scheduler to schedule job onto processor, hence S to P
	pToIOChan := make(chan *Job, 10) // for processor to swap job for IO, hence p to IO. We have 1 IO device, need a buffer for queue

	logger := AuditLogger{SystemTime: clock, Verbose: verbose}

	io := NewIOStream(inputs, ioToSChan, pToIOChan, sToIOSignal, ioCompletedSignal, &logger, clockToIOSignal, clock)
	go io.Run(timeoutCtx, cancel)

	MLFQCfg.Logger = &logger
	MLFQCfg.SToPChan = sToPChan
	MLFQCfg.PToSChan = pToSChan
	MLFQCfg.IOToSChan = ioToSChan
	MLFQCfg.PToSSignal = pToSSignal
	MLFQCfg.SToIOSignal = sToIOSignal

	mlfq := NewMLFQ(*MLFQCfg)
	go mlfq.Run(timeoutCtx)

	processor := NewProcessor(clockToPSignal, sToPChan, pToSSignal, pToSChan, pToIOChan, &logger)
	go processor.Run(timeoutCtx)

	clock.Run(timeoutCtx)

	// requiredCompletion.Wait()

	return logger.SystemOutput
}
