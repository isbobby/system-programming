package main

import (
	"context"
	"fmt"
)

type Processor struct {
	RunningJob  *Job
	SystemClock *Clock
	MLFQ        *MLFQ
	sToPChan    <-chan *Job
	pToIOChan   chan<- *Job
	pToSChan    chan<- *Job
}

func NewProcessor(c *Clock, sToPChan <-chan *Job, pToSChan chan<- *Job, pToIOChan chan<- *Job) Processor {
	return Processor{
		RunningJob:  nil,
		SystemClock: c,
		sToPChan:    sToPChan,
		pToIOChan:   pToIOChan,
		pToSChan:    pToSChan,
	}
}

func (p *Processor) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Time out for simulation")
			return
		default:
		}

		if p.RunningJob != nil {
			// execute and advance time

		} else {
			p.MLFQ.ScheduleJob()

			nextJob, open := <-p.sToPChan
			if !open {
				return
			}

			p.RunningJob = nextJob
		}
	}
}
