package main

import (
	"context"
	"sync/atomic"
	"time"
)

type Clock struct {
	Time          atomic.Int64
	DelayPerCycle time.Duration
	Subscriptions []chan<- interface{}
}

func NewClock(delay time.Duration, subscriptions []chan<- interface{}) *Clock {
	c := Clock{}
	c.Time.Store(0)
	c.DelayPerCycle = delay
	c.Subscriptions = subscriptions
	return &c
}

func (c *Clock) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(c.DelayPerCycle)
			c.publishSignal()
			c.advanceTime()
		}
	}
}

func (c *Clock) advanceTime() {
	c.Time.Add(1)
}

func (c *Clock) publishSignal() {
	for _, subscriber := range c.Subscriptions {
		c.nonBlockingPush(subscriber)
	}
}

func (c *Clock) nonBlockingPush(subscription chan<- interface{}) {
	select {
	case subscription <- struct{}{}:
	default:
	}
}
