package main

import (
	"context"
	"sort"
)

type IOStream struct {
	ScheduledJobs []*Job
	ioToSChan     chan<- *Job
	pToIOChan     <-chan *Job
	SystemTime    Clock
}

func NewIOStream(initialJobs []*Job, ioToSChan chan<- *Job, pToIOChan <-chan *Job) *IOStream {
	return &IOStream{
		ScheduledJobs: initialJobs,
		ioToSChan:     ioToSChan,
		pToIOChan:     pToIOChan,
	}
}

func (s *IOStream) ScheduleInput(ctx context.Context, cancel context.CancelFunc) {
	sort.Slice(s.ScheduledJobs, func(i, j int) bool {
		return s.ScheduledJobs[i].ScheduledTime < s.ScheduledJobs[j].ScheduledTime
	})

	for i := range s.ScheduledJobs {
		job := s.ScheduledJobs[i]

		for job.ScheduledTime > int(s.SystemTime.Time.Load()) {
			break
		}

		s.ioToSChan <- job
	}
}

func (s *IOStream) DoIO(ctx context.Context) {
	for {
		select {
		case job := <-s.pToIOChan:
			IOLog("received job from processor", "ID", job.ID)
			if len(job.InstructionStack) == 0 {
				// TODO log error
				continue
			}

			instruction := job.InstructionStack[len(job.InstructionStack)-1]

			if instruction.IsCPU() {
				s.ioToSChan <- job
			} else if instruction.IsIO() {
				currentSystemTime := s.SystemTime.Time.Load()

				completeTime := int(currentSystemTime) + instruction.Cycle

				IOLog("run IO instruction", "ID", job.ID, "Completion Time", completeTime)

				for completeTime < int(s.SystemTime.Time.Load()) {
					// DO IO
				}

				job.InstructionStack = job.InstructionStack[:len(job.InstructionStack)-1]

				if len(job.InstructionStack) > 0 {
					s.ioToSChan <- job
					IOLog("job IO completed and sent to scheduler", "ID", job.ID)
				} else {
					IOLog("job IO completed", "ID", job.ID)
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
	go s.ScheduleInput(ctx, cancel)
}
