package main

type MLFQ struct {
	MaxPriority        int
	NumQueue           int
	ResetInterval      int
	Queues             []chan Job
	QueueTimeAllotment map[chan Job]int
}

func NewMLFQ(maxPriority int, queueTime []int, resetInterval int, queueSize int) MLFQ {
	timeAllotment := map[chan Job]int{}
	queues := []chan Job{}
	for _, time := range queueTime {
		newQueue := make(chan Job, queueSize)
		timeAllotment[newQueue] = time
		queues = append(queues, newQueue)
	}

	return MLFQ{
		MaxPriority:        maxPriority,
		NumQueue:           len(queues),
		Queues:             queues,
		QueueTimeAllotment: timeAllotment,
		ResetInterval:      resetInterval,
	}
}

func (q *MLFQ) Reset() {

}

func (q *MLFQ) Pop() Job {
	for i := len(q.Queues) - 1; i >= 0; i-- {
		select {
		case job := <-q.Queues[i]:
			return job
		default:
			continue
		}
	}
	return Job{}
}

func (q *MLFQ) Push(j Job) {
	if j.Priority > q.MaxPriority {
		panic("err job with higher priority than MLFQ allows")
	}

	q.Queues[j.Priority] <- j
}
