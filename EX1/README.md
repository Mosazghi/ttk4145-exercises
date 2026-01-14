# Exercise 1

## 1 Sharing a variable

We create two threads and wait for them to finish using `pthread_join(...)`.
We both increment and decrement the shared variable `i` 1 milion times, in both worker functions.
However, `i` doesn't end up with the value `0`. This is due to race condtion that is happening. We
are not using proper synchroization techniques to share `i` correctly, thus ending with an arbitary
final value that may not be consistent.

## 3 Sharing a variable, but properly

We used `pthread_mutex_t` to properly synchroize the two threads. For this purpose we think that it
doesnt really matter which way we go with, but mutex makes more sense in this case.

## Questions

**What is the difference between concurrency and parallelism?**

**concurrency**: task are split into overlapping time periods while **parallelism** run on top of each other (or “at the same time”)

**What is the difference between a race condition and a data race?**

Race _condition_ is a general case, while data race is a specific case of race condition involving simultaneously memory access, at least one being write.

**Very roughly - what does a scheduler do, and how does it do it?**

Decides which decides when and how much time each task gets

**Why would we use multiple threads? What kinds of problems do threads solve?**

Because in a complex program, there are multiple possible events that take place, meaning we have to
keep the UI (e.g.) smooth and responsive while running heavy background tasks, to improve user
experience.

Threads allow tasks to run independently, preventing one slow operation from halting the entire
application, leading to better resource utilization.

**Some languages support “fibers” (sometimes called “green threads”) or “coroutines”? What are they, and why would we rather use them over threads?**

Coroutines are functions that allow execution to be suspended and resumed. Great for cooperative
tasks, exceptions, event loops etc. They are also _stackless_.

Coroutines are managed by the language runtime, while threads are OS level.

**Does creating concurrent programs make the programmer’s life easier? Harder? Maybe both?**

Both. Easier by enabling faster, more responsive applications that utilize modern processors for;
harder due to significant complexity, bugs like race conditions, deadlocks and non-reproducible
bugs, especially with shared mutable state.

**What do you think is best - shared variables or message passing?**

Depends, on the situation. But message passing is more scalable, though slower due to increased
overhead.
