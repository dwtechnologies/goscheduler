package goscheduler

import (
	"fmt"
	"time"
)

const (
	format = "2006-01-02 15:04"
)

type cron struct {
	raw        string
	normalized string
	columns    *cronColumns
	next       *nextDateTime
}

type cronColumns struct {
	minutes     string
	hours       string
	daysOfMonth string
	months      string
	daysOfWeek  string
}

// NextRun will return a time.Time object containing the date and time for when the next execution of the job will be executed.
// Returns time.Time.
func (job *Job) NextRun() time.Time {
	nextRun, _ := time.ParseInLocation(format, fmt.Sprintf("%v-%v-%v %v:%v", job.cron.next.year, *job.cron.convertAndAppendMonth(), *job.cron.convertAndAppendDayOfMonth(), *job.cron.convertAndAppendHour(), *job.cron.convertAndAppendMinute()), time.Local)
	return nextRun
}

// NextRunString will return a string containing the date and time for when the next execution of the job will be executed.
// Returns string.
func (job *Job) NextRunString() string {
	return fmt.Sprintf("%v-%v-%v %v:%v", job.cron.next.year, *job.cron.convertAndAppendMonth(), *job.cron.convertAndAppendDayOfMonth(), *job.cron.convertAndAppendHour(), *job.cron.convertAndAppendMinute())
}

// CronString will return the normalized crontab as a string.
// Returns string.
func (job *Job) CronString() string {
	return job.cron.normalized
}

// CronRawString will return the raw crontab as a string.
// Returns string.
func (job *Job) CronRawString() string {
	return job.cron.raw
}
