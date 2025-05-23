# Multi-level Feedback Queue
## Example output
Scheduled input:
```
inputs := []*Job{
    NewJob(1, 0, []JobInput{IOInstruction}),
    NewJob(2, 0, []JobInput{IOInstruction}),
    NewJob(3, 0, []JobInput{CPUInstruction, CPUInstruction}),
}

MLFQConfig := MLFQConfig{
    QueueConfigs: []QueueConfig{
        {Priority: 4, TimeAllotment: 2},
        {Priority: 3, TimeAllotment: 3},
        {Priority: 2, TimeAllotment: 4},
        {Priority: 1, TimeAllotment: 5},
        {Priority: 0, TimeAllotment: 6},
    },
    ResetInterval: 5,
    QueueSize:     100,
}
```

Observed output logs
```
[T:0][IO] input new job [ID:1]
[T:0][CPU] CPU idle, sent signal for MLFQ 
[T:0][Scheduler] processor idle, control handed to scheduler 
[T:0][Scheduler] MLFQ received job from IO [ID:1]
[T:0][Scheduler] inserting job to queue [Job ID:1 priority:4]
[T:0][Scheduler] Scheduling from highest priority 
[T:0][Scheduler] MLFQ sent job to CPU [ID:1]
[T:0][CPU] swap job out for IO [Job ID:1]
[T:0][CPU] CPU idle, sent signal for MLFQ 
[T:0][IO] received job from processor [ID:1]
[T:0][IO] job has IO instruction, executing [ID:1]
[T:0][IO] run job IO [ID:1 cycle left:5]
[T:0][Scheduler] processor idle, control handed to scheduler 
[T:0][IO] input new job [ID:2]
[T:0][Scheduler] MLFQ received job from IO [ID:2]
[T:0][Scheduler] inserting job to queue [Job ID:2 priority:4]
[T:0][Scheduler] Scheduling from highest priority 
[T:0][Scheduler] MLFQ sent job to CPU [ID:2]
[T:0][CPU] swap job out for IO [Job ID:2]
[T:0][IO] input new job [ID:3]
[T:0][CPU] CPU idle, sent signal for MLFQ 
[T:0][Scheduler] processor idle, control handed to scheduler 
[T:0][Scheduler] MLFQ received job from IO [ID:3]
[T:0][Scheduler] inserting job to queue [Job ID:3 priority:4]
[T:0][Scheduler] Scheduling from highest priority 
[T:0][Scheduler] MLFQ sent job to CPU [ID:3]
[T:0][IO] All jobs scheduled 
[T:1][IO] run job IO [ID:1 cycle left:4]
[T:1][CPU] executing job instruction [Job ID:3 Instruction Left:1 Time Left:2]
[T:2][IO] run job IO [ID:1 cycle left:3]
[T:2][CPU] executing job instruction [Job ID:3 Instruction Left:0 Time Left:1]
[T:2][CPU] complete [Job ID:3]
[T:2][CPU] CPU idle, sent signal for MLFQ 
[T:2][Scheduler] processor idle, control handed to scheduler 
[T:3][IO] run job IO [ID:1 cycle left:2]
[T:4][IO] run job IO [ID:1 cycle left:1]
[T:5][IO] job IO completed [ID:1]
[T:5][IO] received job from processor [ID:2]
[T:5][IO] job has IO instruction, executing [ID:2]
[T:5][IO] run job IO [ID:2 cycle left:5]
[T:6][IO] run job IO [ID:2 cycle left:4]
[T:7][IO] run job IO [ID:2 cycle left:3]
[T:8][IO] run job IO [ID:2 cycle left:2]
[T:9][IO] run job IO [ID:2 cycle left:1]
[T:10][IO] job IO completed [ID:2]
[T:12][Scheduler] scheduler idled for 2 cycles, assume no more tasks, exiting scheduler
```
## Fixes & Change logs
1 - refactored design, where if CPU is blocked on empty task, it will hand over control to MLFQ to schedule the next, simulating handing back control to scheduler

2 - process times out and deadlocks when only job is at IO, IO device does not advance system time

For example, no task on CPU, but a task is in IO being executed, because IO is not advancing time, the task is deadlocked

The fix is the introduction of a logical clock which sends a clock signal to all components.

3 - scheduler termination when no more jobs


When the processor has no more jobs, it hands control over the scheduler

When IO has no more task, it will set its `atomic.bool Busy` to false.

When scheduler has control, sees no pending task from processor, and IO is not busy, and IO has no more tasks to be scheduled for more than 2 clock cycles, it will exit.


4 - priority reshuffle

MLFQ accepts one more clock signal, and when it exceeds a set interval from config, it will perform priority reset to move all jobs into the highest priority.

## To fix
1 - clean up audit logs

2 - use JSON for system configs

3 - refactor in various places

4 - code design summary

