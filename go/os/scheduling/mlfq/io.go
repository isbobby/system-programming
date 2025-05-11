package main

import (
	"context"
	"sort"
	"sync/atomic"
)

type IODevice struct {
	scheduledJobs []*Job

	clockSignal <-chan interface{}
	ioToSChan   chan<- *Job
	pToIOChan   <-chan *Job

	systemTime *Clock

	logger *AuditLogger

	ioDeviceBusy                atomic.Bool
	ioDeviceMoreTasksInSchedule atomic.Bool
}

func NewIODevice(initialJobs []*Job, ioToSChan chan<- *Job, pToIOChan <-chan *Job, logger *AuditLogger, clockSignal <-chan interface{}, clock *Clock) *IODevice {
	return &IODevice{
		scheduledJobs: initialJobs,
		ioToSChan:     ioToSChan,
		pToIOChan:     pToIOChan,
		logger:        logger,
		clockSignal:   clockSignal,
		systemTime:    clock,
	}
}

func (s *IODevice) ScheduleInput(ctx context.Context) {
	s.ioDeviceMoreTasksInSchedule.Store(true)

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

	s.ioDeviceMoreTasksInSchedule.Store(false)
	s.logger.IOLog("All jobs scheduled")
}

func (s *IODevice) PollIOQueue(ctx context.Context) {
	for {
		select {
		case job := <-s.pToIOChan:
			s.doIO(job)
		case <-ctx.Done():
			return
		default:
			continue
		}
	}
}

func (s *IODevice) doIO(job *Job) {
	s.logger.IOLog("received job from processor", "ID", job.ID)
	if len(job.InstructionStack) == 0 {
		s.logger.IOLog("received job from processor", "ID", job.ID)
		return
	}

	instruction := job.InstructionStack[len(job.InstructionStack)-1]

	if instruction.IsCPU() {
		s.logger.IOLog("job has CPU instruction, send back to scheduler", "ID", job.ID)
		s.ioToSChan <- job
		return
	}

	if instruction.IsIO() {
		s.logger.IOLog("job has IO instruction, executing", "ID", job.ID)

		cyclesRemaining := instruction.Cycle

		s.ioDeviceBusy.Store(true)

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

		s.ioDeviceBusy.Store(false)
	}
}

func (s *IODevice) Run(ctx context.Context, cancel context.CancelFunc) {
	go s.PollIOQueue(ctx)
	go s.ScheduleInput(ctx)
}

func (s *IODevice) DeviceBusy() bool {
	return s.ioDeviceBusy.Load()
}

func (s *IODevice) DeviceHasTasks() bool {
	return s.ioDeviceMoreTasksInSchedule.Load()
}

type IODeviceAPI interface {
	DeviceBusy() bool
	DeviceHasTasks() bool
}
