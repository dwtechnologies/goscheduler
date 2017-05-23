package goscheduler

// Scheduler contains the jobs that the scheduler should run as well as a bool for if we should Run or not.
type Scheduler struct {
	run   bool
	count int
	jobs  *[]*Job
}

// Create creates a new Scheduler. Returns *Scheduler.
func Create(f interface{}, p ...interface{}) *Scheduler {
	sched := &Scheduler{
		run:  true,
		jobs: new([]*Job),
	}
	return sched
}

// Start will start the Scheduler and run jobs according to their cron schedule.
func (sched *Scheduler) Start() {
	sched.run = true
}

// Stop will stop all the running jobs.
func (sched *Scheduler) Stop() {
	sched.run = false
}
