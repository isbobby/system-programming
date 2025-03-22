package main

import "fmt"

func main() {
	maxPriority := 3

	mlfqResetInterval := 10
	mlfqQueueSize := 100
	mlfq := NewMLFQ(
		maxPriority,
		[]int{1, 2, 3, 4},
		mlfqResetInterval,
		mlfqQueueSize,
	)

	jobs := []Job{
		NewJob(maxPriority, []JobInput{CPUInstruction}),
	}

	for _, job := range jobs {
		mlfq.Push(job)
	}

	j := mlfq.Pop()
	fmt.Println(j)
}
