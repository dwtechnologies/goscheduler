package goscheduler

import (
	"fmt"
	"reflect"
	"time"
)

// Job contains all the information about a job that can be scheduled.
type Job struct {
	name           *string
	enabled        *bool
	running        *bool
	disableOnError *bool
	count          *int
	err            *error
	cron           *cron
	function       *reflect.Value
	parameters     *[]reflect.Value
	lastResult     *[]reflect.Value
	lastRunTime    *time.Time
	channel        *<-chan time.Time
}

// AddJob takes name, a standard POSIX cron syntax, function and function parameters and turns it into a jobb.
// Returns *Job and error.
func (sched *Scheduler) AddJob(name string, raw string, f interface{}, p ...interface{}) (*Job, error) {
	enabled := true
	cron := &cron{raw: &raw}
	job := &Job{name: &name, enabled: &enabled, cron: cron}

	// Generate a normalized string and check for errors.
	normalized, err := job.cron.normalizeCron()
	if err != nil {
		return nil, err
	}
	job.cron.normalized = normalized

	// Check that the function and parameters are correct.
	function, parameters, err := job.checkFunc(&f, &p)
	if err != nil {
		return nil, err
	}

	job.function = function
	job.parameters = parameters
	*sched.jobs = append(*sched.jobs, job)
	return job, nil
}

// LastRunTime ff
func (job *Job) LastRunTime() (*time.Time, error) {
	if job.lastRunTime == nil {
		return nil, fmt.Errorf("Last Run Time was nil, so it seems it has never been run. In job %v", *job.name)
	}
	return job.lastRunTime, nil
}

// Start will start the specified according to it's cron scheduled.
func (job *Job) Start() {
	go job.run()
}

func (job *Job) run() {
	for t := range *job.channel {
		// Set last run time.
		job.lastRunTime = &t

		// Run the function with the supplied parameters.
		result := job.function.Call(*job.parameters)

		// Check to see if the last result value is error.
		for _, val := range result {
			if val.Interface().(error) != nil {
				err := val.Interface().(error)
				job.err = &err
			}
		}

		// Set the result.
		job.lastResult = &result
	}
}
