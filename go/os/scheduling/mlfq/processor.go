package main

import (
	"context"
)

type Processor struct {
	runningJob  *Job
	SystemClock *Clock
	MLFQ        *MLFQ
	sToPChan    <-chan *Job
	pToIOChan   chan<- *Job
	pToSChan    chan<- *Job
	logger      *Logger
}

func NewProcessor(c *Clock, sToPChan <-chan *Job, pToSChan chan<- *Job, pToIOChan chan<- *Job, logger *Logger) Processor {
	return Processor{
		runningJob:  nil,
		SystemClock: c,
		sToPChan:    sToPChan,
		pToIOChan:   pToIOChan,
		pToSChan:    pToSChan,
		logger:      logger,
	}
}

func (p *Processor) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			p.logger.CPUErrLog("time out signal, processor exiting")
			return
		default:
		}

		if p.runningJob != nil {
			p.runCurrentJob()
		} else {
			nextJob, open := <-p.sToPChan
			if !open {
				p.logger.CPUWarnLog("No more jobs from scheduler, processor exits")
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
		if p.runningJob.TimeAllotment.Load() == 0 {
			p.logger.CPULog("expire job with no more time allotment", "ID", p.runningJob.ID)
			break
		}

		if len(p.runningJob.InstructionStack) == 0 {
			p.logger.CPULog("job with no more instruction completes", "ID", p.runningJob.ID)
			break
		}

		instruction := p.runningJob.InstructionStack[len(p.runningJob.InstructionStack)-1]
		if instruction.IsCPU() {
			p.runningJob.InstructionStack = p.runningJob.InstructionStack[:len(p.runningJob.InstructionStack)-1]
			p.logger.CPULog("run job instruction", "ID", p.runningJob.ID, "Instruction Left", len(p.runningJob.InstructionStack), "Time Left", p.runningJob.TimeAllotment.Load())
			p.runningJob.TimeAllotment.Add(-1)
			p.SystemClock.AdvanceTime()
		} else {
			p.pToIOChan <- p.runningJob
			p.logger.CPULog("sent job to IO device", "ID", p.runningJob.ID)
			p.runningJob = nil
			return
		}
	}

	if len(p.runningJob.InstructionStack) > 0 {
		p.logger.CPULog("sent expired job to MLFQ", "ID", p.runningJob.ID)
		p.pToSChan <- p.runningJob
	}

	p.runningJob = nil
}
