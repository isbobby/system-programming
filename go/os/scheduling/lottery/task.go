package main

type Schedulable interface {
	ID() int
}

type task struct {
	id          int
	ticketCount int
}

func NewTask(id int, ticketCount int) Schedulable {
	return &task{
		id:          id,
		ticketCount: ticketCount,
	}
}

func (t *task) ID() int {
	return t.id
}
