package main

import (
	"context"
	"fmt"
)

type Processor struct {
	runningJob  *Job
	SystemClock *Clock
	MLFQ        *MLFQ
	sToPChan    <-chan *Job
	pToIOChan   chan<- *Job
	pToSChan    chan<- *Job
}

func NewProcessor(c *Clock, sToPChan <-chan *Job, pToSChan chan<- *Job, pToIOChan chan<- *Job) Processor {
	return Processor{
		runningJob:  nil,
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

		if p.runningJob != nil {
			p.runCurrentJob()
		} else {
			nextJob, open := <-p.sToPChan
			if !open {
				fmt.Println("No more jobs from scheduler, processor exits")
				return
			}

			p.runningJob = nextJob
		}
	}
}

func (p *Processor) runCurrentJob() {
	if len(p.runningJob.InstructionStack) == 0 {
		p.runningJob = nil
		return
	}

	// CPU Cycle
	for {
		fmt.Println("[SYS TIME]:", p.SystemClock.Time.Load())
		if p.runningJob.TimeAllotment.Load() == 0 {
			fmt.Println("(P) Job no more time allotment, expire", p.runningJob.ID)
			break
		}

		if len(p.runningJob.InstructionStack) == 0 {
			fmt.Println("(P) Job no more instruction, exit")
			break
		}

		instruction := p.runningJob.InstructionStack[len(p.runningJob.InstructionStack)-1]
		fmt.Println("(P) Executing job", p.runningJob.ID, "instruction:", instruction, "Time left", p.runningJob.TimeAllotment.Load())

		if instruction.IsCPU() {
			p.runningJob.InstructionStack = p.runningJob.InstructionStack[:len(p.runningJob.InstructionStack)-1]
			p.runningJob.TimeAllotment.Add(-1)
			p.SystemClock.AdvanceTime()
		} else {
			p.pToIOChan <- p.runningJob
			fmt.Println("(P) Sent IO job", p.runningJob.ID, "back to IO")
			p.runningJob = nil
			return
		}
	}

	if len(p.runningJob.InstructionStack) > 0 {
		fmt.Println("(P) Sending expired job to MLFQ")
		p.pToSChan <- p.runningJob
	}

	p.runningJob = nil
}
