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
	pToSSignal  chan<- interface{}
	logger      *AuditLogger
}

func NewProcessor(c *Clock, sToPChan <-chan *Job, pToSSignal chan<- interface{}, pToSChan chan<- *Job, pToIOChan chan<- *Job, logger *AuditLogger) Processor {
	return Processor{
		runningJob:  nil,
		SystemClock: c,
		sToPChan:    sToPChan,
		pToSSignal:  pToSSignal,
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
			// if non-OS job is scheduled, run it
			p.runCurrentJob()
		} else {
			// otherwise, signal scheduler to run on processor
			p.logger.CPUWarnLog("CPU idle, sent signal for MLFQ")
			p.pToSSignal <- struct{}{}

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
		if len(p.runningJob.InstructionStack) == 0 {
			p.logger.CPUAuditLog(COMPLETE, JobIDKey, p.runningJob.ID)
			break
		}

		if p.runningJob.TimeAllotment.Load() == 0 {
			p.logger.CPUAuditLog(EXPIRE, JobIDKey, p.runningJob.ID)
			break
		}

		instruction := p.runningJob.InstructionStack[len(p.runningJob.InstructionStack)-1]
		if instruction.IsCPU() {
			p.runningJob.InstructionStack = p.runningJob.InstructionStack[:len(p.runningJob.InstructionStack)-1]
			p.logger.CPUAuditLog(EXEC, JobIDKey, p.runningJob.ID, "Instruction Left", len(p.runningJob.InstructionStack), "Time Left", p.runningJob.TimeAllotment.Load())
			p.runningJob.TimeAllotment.Add(-1)
			p.SystemClock.AdvanceTime()
		} else {
			p.pToIOChan <- p.runningJob
			p.logger.CPUAuditLog(SWAP, JobIDKey, p.runningJob.ID)
			p.runningJob = nil
			return
		}
	}

	if len(p.runningJob.InstructionStack) > 0 {
		p.pToSChan <- p.runningJob
	}

	p.runningJob = nil
}
