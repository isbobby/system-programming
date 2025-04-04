package main

import (
	"context"
	"sort"
)

type IOStream struct {
	scheduledJobs     []*Job
	ioCompletedSignal chan<- interface{}

	ioToSChan   chan<- *Job
	pToIOChan   <-chan *Job
	sToIOSignal <-chan interface{}
	clockSignal <-chan interface{}

	systemTime *Clock

	logger *AuditLogger
}

func NewIOStream(initialJobs []*Job, ioToSChan chan<- *Job, pToIOChan <-chan *Job, sToIOSignal <-chan interface{}, ioCompletedSignal chan<- interface{}, logger *AuditLogger, clockSignal <-chan interface{}, clock *Clock) *IOStream {
	return &IOStream{
		scheduledJobs:     initialJobs,
		ioToSChan:         ioToSChan,
		pToIOChan:         pToIOChan,
		sToIOSignal:       sToIOSignal,
		ioCompletedSignal: ioCompletedSignal,
		logger:            logger,
		clockSignal:       clockSignal,
		systemTime:        clock,
	}
}

func (s *IOStream) ScheduleInput(ctx context.Context) {
	sort.Slice(s.scheduledJobs, func(i, j int) bool {
		return s.scheduledJobs[i].ScheduledTime < s.scheduledJobs[j].ScheduledTime
	})

	for len(s.scheduledJobs) > 0 {
		job := s.scheduledJobs[0]
		s.scheduledJobs = s.scheduledJobs[1:]

		for job.ScheduledTime > int(s.systemTime.Time.Load()) {
		}

		s.logger.IOLog("input new job", "ID", job.ID)
		s.ioToSChan <- job
	}

	s.logger.IOLog("All jobs scheduled")
	s.ioCompletedSignal <- struct{}{}
}

func (s *IOStream) DoIO(ctx context.Context) {
	for {
		select {
		case job := <-s.pToIOChan:
			s.logger.IOLog("received job from processor", "ID", job.ID)
			if len(job.InstructionStack) == 0 {
				// TODO log error
				continue
			}

			instruction := job.InstructionStack[len(job.InstructionStack)-1]

			if instruction.IsCPU() {
				s.logger.IOLog("job has CPU instruction, send back to scheduler", "ID", job.ID)
				s.ioToSChan <- job
			} else if instruction.IsIO() {
				s.logger.IOLog("job has IO instruction, executing", "ID", job.ID)

				cyclesRemaining := instruction.Cycle

				for cyclesRemaining > 0 {
					// note what if there's existing cycle signal
					s.logger.IOLog("run job IO", "ID", job.ID, "cycle left", cyclesRemaining)
					<-s.clockSignal
					cyclesRemaining -= 1
				}

				job.InstructionStack = job.InstructionStack[:len(job.InstructionStack)-1]

				if len(job.InstructionStack) > 0 {
					s.ioToSChan <- job
					s.logger.IOLog("job IO completed and sent to scheduler", "ID", job.ID)
				} else {
					s.logger.IOLog("job IO completed", "ID", job.ID)
				}
			}

		case <-ctx.Done():
			return
		default:
			continue
		}
	}
}

func (s *IOStream) Run(ctx context.Context, cancel context.CancelFunc) {
	go s.DoIO(ctx)
	go s.ScheduleInput(ctx)
}
