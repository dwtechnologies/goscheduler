package goscheduler

import (
	"fmt"
	"reflect"
	"time"
)

// Job contains all the information about a job that can be scheduled.
type Job struct {
	name           string
	enabled        bool
	running        bool
	disableOnError bool
	count          int

	// Global
	global *global

	// Cron and function
	cron       *cron
	function   *reflect.Value
	parameters *[]reflect.Value

	// Result
	result *result

	// Trigger
	timer *time.Timer
}

type result struct {
	result *[]reflect.Value
	time   *time.Time
	err    error
}

// AddJob takes name, a standard POSIX cron syntax, function and function parameters and turns it into a jobb.
// Returns *Job and error.
func (sched *Scheduler) AddJob(name string, raw string, f interface{}, p ...interface{}) (*Job, error) {
	cron := &cron{raw: raw}
	job := &Job{
		name:           name,
		disableOnError: true,
		cron:           cron,
		timer:          nil,
		global:         sched.global,
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

	return job, nil
}

// LastRunTime return the date and time of when the specified job was last run. Returns *time.Time
func (job *Job) LastRunTime() (*time.Time, error) {
	if job.result.time == nil {
		return nil, fmt.Errorf("Last Run Time was nil, so it seems it has never been run. In job %v", job.name)
	}
	return job.result.time, nil
}

// Running will return true if the job is currently running. Otherwise returns false.
// Returns bool.
func (job *Job) Running() bool {
	return job.running
}

// Enabled will return true if the job is currently enabled or not. Otherwise returns false.
// Returns bool.
func (job *Job) Enabled() bool {
	return job.enabled
}

// Stop will stop the specified job from running. It will however let the current running job finish before exiting.
func (job *Job) Stop() {
	job.enabled = false
}

// Start will start the specified according to it's cron scheduled.
func (job *Job) Start() {
	job.enabled = true
	go job.run()
}

// Run the job, should be in a seperate go routine.
func (job *Job) run() {
	for {
		// If local or global enabled has been turned off, exit the go routine.
		switch {
		case !job.global.enabled:
			job.Stop()
			return

		case !job.enabled:
			return
		}

		// Generate the next run.
		job.genTimer()
		for t := range job.timer.C {
			job.running = true
			job.global.count++

			// Run the function and check to see if the last result value is error.
			result := job.function.Call(*job.parameters)
			job.makeResult(&result, &t)

			job.running = false
			job.global.count++
			break
		}
		job.timer.Stop()
	}
}

// makeResult will generate the result and any errors.
func (job *Job) makeResult(r *[]reflect.Value, t *time.Time) {
	result := &result{
		result: r,
		time:   t,
	}
	// Loop through r and see if we have any errors.
	for _, val := range *r {
		if val.Interface().(error) != nil {
			err := val.Interface().(error)
			result.err = err
		}
	}
}

// genTimer will generate the next time for the job to run and set the time.Timer to the diff value between now and the time in the future.
func (job *Job) genTimer() {
	diff := job.genNextRun()
	job.timer = time.NewTimer(time.Duration(int(*diff)) * time.Second)
}
