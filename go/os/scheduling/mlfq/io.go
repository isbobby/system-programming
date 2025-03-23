package main

import "context"

type IOStream struct {
	ScheduledJobs []Job
	ioToSChan     chan<- *Job
	pToIOChan     <-chan *Job
}

func NewIOStream(initialJobs []Job, ioToSChan chan<- *Job, pToIOChan <-chan *Job) *IOStream {
	return &IOStream{
		ScheduledJobs: initialJobs,
		ioToSChan:     ioToSChan,
		pToIOChan:     pToIOChan,
	}
}

func (s *IOStream) ScheduleInput(ctx context.Context, cancel context.CancelFunc) {
	for {

	}
}

func (s *IOStream) DoIO(ctx context.Context, job *Job) {

}
