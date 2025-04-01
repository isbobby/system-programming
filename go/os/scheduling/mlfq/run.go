package main

import (
	"context"
	"time"
)

func RunSystem(MLFQCfg *MLFQConfig, inputs []*Job, timeout time.Duration, verbose bool) []AuditLog {
	ctx := context.Background()
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	clock := NewClock()
	pToSSignal := make(chan interface{})
	ioToSChan := make(chan *Job, 10) // for IO to enqueue job for scheduling, hence io To S, assuming 10 job buffer
	pToSChan := make(chan *Job, 1)   // for processor to expire a job back to scheduler, hence p to S
	sToPChan := make(chan *Job)      // for scheduler to schedule job onto processor, hence S to P
	pToIOChan := make(chan *Job, 10) // for processor to swap job for IO, hence p to IO. We have 1 IO device, need a buffer for queue

	logger := AuditLogger{SystemTime: clock, Verbose: verbose}

	io := NewIOStream(inputs, ioToSChan, pToIOChan, &logger, clock)
	go io.Run(timeoutCtx, cancel)

	MLFQCfg.Logger = &logger
	MLFQCfg.SToPChan = sToPChan
	MLFQCfg.PToSChan = pToSChan
	MLFQCfg.IOToSChan = ioToSChan
	MLFQCfg.PToSSignal = pToSSignal

	mlfq := NewMLFQ(*MLFQCfg)
	go mlfq.Run(timeoutCtx)

	processor := NewProcessor(clock, sToPChan, pToSSignal, pToSChan, pToIOChan, &logger)
	processor.Run(timeoutCtx)

	return logger.SystemOutput
}
