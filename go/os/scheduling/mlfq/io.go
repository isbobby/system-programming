package main

import (
	"context"
	"fmt"
)

type IOStream struct {
	ScheduledJobs []*Job
	ioToSChan     chan<- *Job
	pToIOChan     <-chan *Job
}

func NewIOStream(initialJobs []*Job, ioToSChan chan<- *Job, pToIOChan <-chan *Job) *IOStream {
	return &IOStream{
		ScheduledJobs: initialJobs,
		ioToSChan:     ioToSChan,
		pToIOChan:     pToIOChan,
	}
}

func (s *IOStream) ScheduleInput(ctx context.Context, cancel context.CancelFunc) {
	for i := range s.ScheduledJobs {
		job := s.ScheduledJobs[i]
		fmt.Println("Scheduled")
		s.ioToSChan <- job
	}
}

func (s *IOStream) DoIO(ctx context.Context, job *Job) {
}

func (s *IOStream) Run(ctx context.Context, cancel context.CancelFunc) {
	go s.ScheduleInput(ctx, cancel)
}
