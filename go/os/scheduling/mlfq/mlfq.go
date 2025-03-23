package main

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
}

func NewMLFQ(maxPriority int, queueTime []int, resetInterval int, queueSize int, clock *Clock, sToPChan chan<- *Job, ioToSChan <-chan *Job, pToSChan <-chan *Job) MLFQ {
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

func (q *MLFQ) ScheduleJob() {
	for i := q.MaxPriority; i >= 0; i-- {
		select {
		case job := <-q.Queues[i]:
			q.sToPChan <- job
		default:
			continue
		}
	}
	close(q.sToPChan)
}

func (q *MLFQ) Push(j *Job) {
	if j.Priority > q.MaxPriority {
		panic("err job with higher priority than MLFQ allows")
	}

	timeAlloted := q.QueueTimeAllotment[j.Priority]

	j.TimeAllotment.Store(int32(timeAlloted))

	q.Queues[j.Priority] <- j
}
