package main

import (
	"context"
	"sync"
	"time"
)

func RunSystem(MLFQCfg *MLFQConfig, inputs []*Job, timeout time.Duration, verbose bool) []AuditLog {
	ctx := context.Background()
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	clockToIOSignal := make(chan interface{})
	clockToSSignal := make(chan interface{})
	clockToPSignal := make(chan interface{})
	clockToResetSSignal := make(chan interface{})
	clockSubscriptions := []chan<- interface{}{clockToIOSignal, clockToPSignal, clockToSSignal, clockToResetSSignal}
	clock := NewClock(time.Duration(100*time.Millisecond), clockSubscriptions)

	pToSSignal := make(chan interface{})

	ioToSChan := make(chan *Job)     // for IO to enqueue job for scheduling, hence io To S, assuming 10 job buffer
	pToIOChan := make(chan *Job, 10) // for processor to swap job for IO, hence p to IO. We have 1 IO device, need a buffer for queue
	pToSChan := make(chan *Job, 1)   // for processor to expire a job back to scheduler, hence p to S
	sToPChan := make(chan *Job)      // for scheduler to schedule job onto processor, hence S to P

	logger := AuditLogger{SystemTime: clock, Verbose: verbose}

	io := NewIODevice(inputs, ioToSChan, pToIOChan, &logger, clockToIOSignal, clock)
	go io.Run(timeoutCtx, cancel)

	MLFQCfg.Logger = &logger
	MLFQCfg.SToPChan = sToPChan
	MLFQCfg.PToSChan = pToSChan
	MLFQCfg.IOToSChan = ioToSChan
	MLFQCfg.PToSSignal = pToSSignal
	MLFQCfg.ClockSignal = clockToSSignal
	MLFQCfg.ResetClockSignal = clockToResetSSignal

	mlfq := NewMLFQ(*MLFQCfg, io)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		mlfq.Run(timeoutCtx)
		wg.Done()
	}()

	processor := NewProcessor(clockToPSignal, sToPChan, pToSSignal, pToSChan, pToIOChan, &logger)
	go processor.Run(timeoutCtx)

	go clock.Run(timeoutCtx)

	wg.Wait()

	return logger.SystemOutput
}
