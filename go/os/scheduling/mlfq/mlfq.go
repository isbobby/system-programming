package main

import (
	"context"
)

type MLFQ struct {
	MaxPriority        int
	NumQueue           int
	ResetInterval      int
	Queues             []chan *Job
	QueueTimeAllotment map[int]int

	sToPChan chan<- *Job

	pToSChan  <-chan *Job
	ioToSChan <-chan *Job

	SystemClock *Clock

	logger *Logger
}

func NewMLFQ(maxPriority int, queueTime []int, resetInterval int, queueSize int, clock *Clock, sToPChan chan<- *Job, ioToSChan <-chan *Job, pToSChan <-chan *Job, logger *Logger) MLFQ {
	timeAllotment := map[int]int{}
	queues := []chan *Job{}
	for priority, time := range queueTime {
		newQueue := make(chan *Job, queueSize)
		timeAllotment[priority] = time
		queues = append(queues, newQueue)
	}

	return MLFQ{
		MaxPriority:        maxPriority,
		NumQueue:           len(queues),
		Queues:             queues,
		QueueTimeAllotment: timeAllotment,
		ResetInterval:      resetInterval,

		sToPChan:  sToPChan,
		ioToSChan: ioToSChan,
		pToSChan:  pToSChan,
		logger:    logger,
	}
}

func (q *MLFQ) Reset() {
	jobs := []*Job{}

	// remove all jobs from non-max priority queues
	for i := q.MaxPriority - 1; i >= 0; i-- {
		select {
		case job := <-q.Queues[i]:
			jobs = append(jobs, job)
		default:
			continue
		}
	}

	for _, job := range jobs {
		q.Queues[q.MaxPriority] <- job
	}
}

func (q *MLFQ) AcceptJobFromIO(ctx context.Context) {
	for {
		select {
		case job := <-q.ioToSChan:
			q.logger.MLFQLog("MLFQ received job from IO", "ID", job.ID)
			q.push(job)
		case <-ctx.Done():
			return
		default:
			continue
		}
	}
}

func (q *MLFQ) AcceptExpiredJobFromProc(ctx context.Context) {
	for {
		select {
		case job := <-q.pToSChan:
			q.logger.MLFQLog("MLFQ received expired job from CPU", "ID", job.ID)
			job.DecreasePriority()
			q.push(job)
		case <-ctx.Done():
			return
		default:
			continue
		}
	}
}

func (q *MLFQ) ScheduleJob(ctx context.Context) {
	for {
		for i := q.MaxPriority; i >= 0; i-- {
			select {
			case job := <-q.Queues[i]:
				q.sToPChan <- job
				q.logger.MLFQLog("MLFQ sent job to CPU", "ID", job.ID)
			case <-ctx.Done():
				close(q.sToPChan)
				return
			default:
				continue
			}
		}
	}
}

func (q *MLFQ) push(j *Job) {
	if j.Priority == nil {
		j.Priority = &q.MaxPriority
	}

	if *j.Priority > q.MaxPriority {
		panic("err job with higher priority than MLFQ allows")
	}

	timeAlloted := q.QueueTimeAllotment[*j.Priority]

	j.TimeAllotment.Store(int32(timeAlloted))

	q.Queues[*j.Priority] <- j
}

func (q *MLFQ) Run(ctx context.Context) {
	go q.AcceptExpiredJobFromProc(ctx)

	go q.AcceptJobFromIO(ctx)

	go q.ScheduleJob(ctx)

	go q.Reset()
}
