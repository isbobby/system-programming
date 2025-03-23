package main

import "sync/atomic"

type Clock struct {
	Time atomic.Int64
}

func NewClock() *Clock {
	c := Clock{}
	c.Time.Store(0)
	return &c
}

func (c *Clock) AdvanceTime() {
	c.Time.Add(1)
}
