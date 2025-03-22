# Scheduler Implementation
Following is an attempt to implement scheduler simulation with three components
1. Input Stream - represents IO
2. Scheduler - maintains a buffer of ready tasks, apply scheduling strategy to choose a task, and schedule it for execution
3. Processor - wait for scheduler to schedule and execute tasks

Some other important constructs are
1. `ctx context.WithTimeout()` the same context is propagated across different go routines to help identify deadlocks.
2. `var systemTime atomic.Int32` is used to synchronise task timing across the routines
3. `var ExecStats` help with logging and visualisation.

## Code Design
![](scheduler_code_fifo.png)

## Todo
Current implementation needs the following improvement
1. Calculate response time in addition to turnaround time
2. Verify if the existing architecture can support task switching
