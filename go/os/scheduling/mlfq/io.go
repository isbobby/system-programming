package main

import (
	"context"
	"fmt"
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
			fmt.Println("(IO) received job from P to do IO", job.ID)
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

				fmt.Println("(IO) Doing IO, expected complet time:", completeTime)

				for completeTime < int(s.SystemTime.Time.Load()) {
					// DO IO
				}

				job.InstructionStack = job.InstructionStack[:len(job.InstructionStack)-1]

				if len(job.InstructionStack) > 0 {
					// back to S without changing priority
					fmt.Println("(IO) job IO completed, ID:", job.ID, "sending back to scheduler")
					s.ioToSChan <- job
				} else {
					fmt.Println("(IO) job IO completed, ID:", job.ID)
					// Log job completion
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
