package main

type TaskQueue interface {
	AddTask(Schedulable) error
	RemoveTask(id int) (Schedulable, error)
	FindTask(ticket int) (Schedulable, error)
	Tasks() []Schedulable
}
