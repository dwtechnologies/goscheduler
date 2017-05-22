package goscheduler

import (
	"fmt"
	"time"
)

// Scheduler contains the jobs that the scheduler should run as well as a bool for if we should Run or not.
type Scheduler struct {
	running *bool
	ticker  *time.Ticker
	jobs    *[]*Job
}

// Create creates a new Scheduler. Returns *Scheduler.
func Create(f interface{}, p ...interface{}) *Scheduler {
	sched := &Scheduler{jobs: new([]*Job)}
	return sched
}

// Start will start the Scheduler and run jobs according to their cron schedule.
func (sched *Scheduler) Start() {
	sched.ticker = time.NewTicker(time.Second * 1)
	go sched.run()

	running := true
	sched.running = &running
}

func (sched *Scheduler) run() {
	for t := range sched.ticker.C {
		fmt.Println("hejsan", t)
	}
}

// Stop will stop the Scheduler.
func (sched *Scheduler) Stop() {
	sched.ticker.Stop()

	running := false
	sched.running = &running
}
