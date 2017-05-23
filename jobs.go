package goscheduler

import (
	"fmt"
	"reflect"
	"time"
)

// Job contains all the information about a job that can be scheduled.
type Job struct {
	name    string
	enabled bool

	running        bool
	disableOnError bool
	quit           bool
	count          int
	err            error

	// Cron and function
	cron        *cron
	function    *reflect.Value
	parameters  *[]reflect.Value
	lastResult  *[]reflect.Value
	lastRunTime *time.Time

	// Trigger
	timer *time.Timer
}

// AddJob takes name, a standard POSIX cron syntax, function and function parameters and turns it into a jobb.
// Returns *Job and error.
func (sched *Scheduler) AddJob(name string, raw string, f interface{}, p ...interface{}) (*Job, error) {
	cron := &cron{raw: raw}
	job := &Job{
		name:    name,
		enabled: true,
		cron:    cron,
		timer:   nil,
	}

	// Generate a normalized string and check for errors.
	normalized, err := job.cron.normalizeCron()
	if err != nil {
		return nil, err
	}
	job.cron.normalized = *normalized

	// Check that the function and parameters are correct.
	function, parameters, err := job.checkFunc(&f, &p)
	if err != nil {
		return nil, err
	}

	job.function = function
	job.parameters = parameters
	*sched.jobs = append(*sched.jobs, job)

	// Generate next execution time and add it to the timer.
	job.Start()

	return job, nil
}

// LastRunTime return the date and time of when the specified job was last run. Returns *time.Time
func (job *Job) LastRunTime() (*time.Time, error) {
	if job.lastRunTime == nil {
		return nil, fmt.Errorf("Last Run Time was nil, so it seems it has never been run. In job %v", job.name)
	}
	return job.lastRunTime, nil
}

// Stop will stop the specified job from running. It will however let the current running job finish before exiting.
func (job *Job) Stop() {
	job.quit = true
}

// Start will start the specified according to it's cron scheduled.
func (job *Job) Start() {
	go job.run()
}

func (job *Job) run() {
	for {
		job.genTimer()
		for t := range job.timer.C {
			job.lastRunTime = &t
			result := job.function.Call(*job.parameters)

			// Check to see if the last result value is error.
			job.checkFuncError(&result)

			job.lastResult = &result
			break
		}
		job.timer.Stop()

		if job.quit {
			return
		}
	}
}

func (job *Job) checkFuncError(result *[]reflect.Value) {
	for _, val := range *result {
		if val.Interface().(error) != nil {
			err := val.Interface().(error)
			job.err = err
		}
	}
}

func (job *Job) genTimer() {
	diff := job.genNextRun()
	job.timer = time.NewTimer(time.Duration(int(*diff)) * time.Second)
}
