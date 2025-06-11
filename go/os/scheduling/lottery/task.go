package main

type Schedulable interface {
	ID() int
	MatchInterval(int) bool
	Interval() [2]int
	SetInterval([2]int)
}

type task struct {
	ticketInterval [2]int
	id             int
	ticketCount    int
}

func NewSchedulableTask(id int, ticketCount int, ticketInterval [2]int) Schedulable {
	return &task{
		id:             id,
		ticketCount:    ticketCount,
		ticketInterval: ticketInterval,
	}
}

func (t *task) ID() int {
	return t.id
}

func (t *task) MatchInterval(ticket int) bool {
	return t.ticketInterval[0] <= ticket && ticket <= t.ticketInterval[1]
}

func (t *task) Interval() [2]int {
	return t.ticketInterval
}

func (t *task) SetInterval(newTicketInterval [2]int) {
	t.ticketInterval = newTicketInterval
}
