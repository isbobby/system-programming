# Multi-level Feedback Queue

## Features
1 - refactored design, where if CPU is blocked on empty task, it will hand over control to MLFQ to schedule the next, simulating handing back control to scheduler
2 - code time out and deadlocks when only job is at IO, IO device does not advance system time
    1.1 - for example, no task on CPU, but a task is in IO being executed, because IO is not advancing time, the task is deadlocked
    1.2 - test case: one task with IO instruction only
3 - scheduler termination implementation
    3.1 - when processor has no more task, it hands control over the scheduler
    3.2 - when IO has no more task, it will set its `atomic.bool Busy` to false
    3.3 - when scheduler has control, sees no pending task from processor, and IO is not busy, and IO has no more tasks to be scheduled it will exit
4 - priority reshuffle
    4.1 MLFQ accepts a clock signal, and when it exceeds a set interval from config, it will perform priority reset to move all jobs into the highest priority

## To fix
1 - clean up audit logs
2 - use JSON for system configs
3 - refactor in various places
4 - code design summary

