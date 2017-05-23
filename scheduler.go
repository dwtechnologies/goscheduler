package goscheduler

import "time"

// Scheduler defines the scheduler main construct.
type Scheduler struct {
	jobs   *[]*Job
	count  int
	global *global
}

type global struct {
	enabled bool
	count   int
}

// Create creates a new Scheduler.
// Returns *Scheduler.
func Create() *Scheduler {
	sched := &Scheduler{
		jobs:   new([]*Job),
		global: new(global),
	}
	return sched
}

// Start will start the Scheduler and run jobs according to their cron schedule.
func (sched *Scheduler) Start() {
	sched.global.enabled = true
}

// Stop will stop all the running jobs.
func (sched *Scheduler) Stop() {
	sched.global.enabled = true
}

// Status will return the current status of the Scheduler.
// Returns true on enabled. False on disabled.
func (sched *Scheduler) Status() bool {
	return sched.global.enabled
}

// Done will return true when all Jobs have finished after Scheduler sent Stop().
// Returns bool.
func (sched *Scheduler) Done() bool {
	if sched.global.count < 1 {
		return true
	}
	return false
}

// WaitDone will wait until Done() is false and return. Checks every 100ms.
func (sched *Scheduler) WaitDone() {
	for {
		if sched.Done() {
			break
		}

		// Check if jobs have finished every second.
		time.Sleep(time.Duration(100) * time.Millisecond)
	}
}
