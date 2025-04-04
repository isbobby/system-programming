package main

import (
	"context"
	"errors"
	"fmt"
	"sort"
)

type QueueConfig struct {
	Priority      int
	TimeAllotment int
}

type MLFQConfig struct {
	QueueConfigs  []QueueConfig
	ResetInterval int
	QueueSize     int

	SToPChan    chan<- *Job
	IOToSChan   <-chan *Job
	PToSChan    <-chan *Job
	PToSSignal  <-chan interface{}
	SToIOSignal chan<- interface{}
	IOToSSignal <-chan interface{}

	Logger *AuditLogger
}

func (cfg MLFQConfig) Validate() error {
	if cfg.SToPChan == nil || cfg.IOToSChan == nil || cfg.PToSChan == nil || cfg.PToSSignal == nil || cfg.SToIOSignal == nil {
		return errors.New("err attempt to initialise MLFQ with some nil channel")
	}

	if cfg.QueueSize < 1 {
		return errors.New("err MLFQ queue size must be at least 1")
	}

	sort.Slice(cfg.QueueConfigs, func(a, b int) bool {
		return cfg.QueueConfigs[a].Priority < cfg.QueueConfigs[b].Priority
	})

	for i := len(cfg.QueueConfigs) - 1; i >= 0; i-- {
		if cfg.QueueConfigs[i].Priority != i {
			return fmt.Errorf("err priority level, expected %v, got %v", i, cfg.QueueConfigs[i].Priority)
		}
	}

	return nil
}

type MLFQ struct {
	MaxPriority        int
	NumQueue           int
	ResetInterval      int
	QueuesByPriority   map[int]chan *Job
	QueueTimeAllotment map[int]int

	sToPChan    chan<- *Job
	pToSChan    <-chan *Job
	pToSSignal  <-chan interface{}
	ioToSChan   <-chan *Job
	sToIOSignal chan<- interface{}
	ioToSSignal <-chan interface{}

	logger *AuditLogger
}

func NewMLFQ(cfg MLFQConfig) MLFQ {
	if err := cfg.Validate(); err != nil {
		panic(fmt.Errorf("failed to initialise MLFQ due to config error: %v", err))
	}

	mlfq := MLFQ{
		ResetInterval: cfg.ResetInterval,

		pToSSignal:  cfg.PToSSignal,
		sToIOSignal: cfg.SToIOSignal,
		ioToSSignal: cfg.IOToSSignal,

		sToPChan:  cfg.SToPChan,
		ioToSChan: cfg.IOToSChan,
		pToSChan:  cfg.PToSChan,
		logger:    cfg.Logger,
	}

	timeAllotment := map[int]int{}
	QueuesByPriority := map[int]chan *Job{}
	for _, config := range cfg.QueueConfigs {
		newQueue := make(chan *Job, cfg.QueueSize)
		timeAllotment[config.Priority] = config.TimeAllotment
		QueuesByPriority[config.Priority] = newQueue
		mlfq.MaxPriority = max(mlfq.MaxPriority, config.Priority)
	}

	mlfq.QueueTimeAllotment = timeAllotment
	mlfq.QueuesByPriority = QueuesByPriority

	return mlfq
}

func (q *MLFQ) Reset() {
	jobs := []*Job{}

	// remove all jobs from non-max priority queues, making use of reset interval
	for prio, queue := range q.QueuesByPriority {
		if prio == q.MaxPriority {
			continue
		}

		for len(queue) > 0 {
			jobs = append(jobs, <-queue)
		}
	}

	for _, job := range jobs {
		q.QueuesByPriority[q.MaxPriority] <- job
	}
}

func (q *MLFQ) AcceptJobFromIO(ctx context.Context) {
	select {
	case job := <-q.ioToSChan:
		q.logger.MLFQLog("MLFQ received job from IO", "ID", job.ID)
		q.push(job)
	case <-ctx.Done():
		return
	default:
		return
	}
}

func (q *MLFQ) AcceptExpiredJobFromProc(ctx context.Context) {
	select {
	case job := <-q.pToSChan:
		job.DecreasePriority()
		q.logger.MLFQLog("MLFQ received expired job from CPU", "ID", job.ID, "New Priority", *job.Priority)
		q.push(job)
	case <-ctx.Done():
		return
	default:
		return
	}
}

func (q *MLFQ) push(j *Job) {
	if j.Priority == nil {
		var newJobPriority int = q.MaxPriority
		j.Priority = &newJobPriority
	}

	q.logger.MLFQLog("inserting job to queue", JobIDKey, j.ID, "priority", *j.Priority)

	if *j.Priority > q.MaxPriority {
		panic("err job with higher priority than MLFQ allows")
	}

	timeAlloted := q.QueueTimeAllotment[*j.Priority]

	j.TimeAllotment.Store(int32(timeAlloted))

	q.QueuesByPriority[*j.Priority] <- j
}

func (q *MLFQ) pollForReadyTasks(ctx context.Context) bool {
	for {
		select {
		case <-ctx.Done():
			return true
		default:
			for i := range q.QueuesByPriority {
				if len(q.QueuesByPriority[i]) > 0 {
					return true
				}
			}
			return false
		}
	}
}

func (q *MLFQ) HandleProcSignal(ctx context.Context) {
	for {
		<-q.pToSSignal
		q.logger.MLFQLog("processor idle, control handed to scheduler")
		for {
			// if no ready tasks, attempt to receive from IO/CPU
			if ready := q.pollForReadyTasks(ctx); !ready {
				q.AcceptJobFromIO(ctx)
				q.AcceptExpiredJobFromProc(ctx)
			} else {
				break
			}
		}
		q.scheduleJob(ctx)
	}
}

func (q *MLFQ) scheduleJob(ctx context.Context) {
	q.logger.MLFQLog("Scheduling from highest priority")
	for i := q.MaxPriority; i >= 0; i-- {
		select {
		case job := <-q.QueuesByPriority[i]:
			q.sToPChan <- job
			q.logger.MLFQLog("MLFQ sent job to CPU", "ID", job.ID)
			return
		case <-ctx.Done():
			close(q.sToPChan)
			return
		default:
			continue
		}
	}
}

func (q *MLFQ) Run(ctx context.Context) {
	go q.HandleProcSignal(ctx)
}
