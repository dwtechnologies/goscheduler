package goscheduler

type cron struct {
	raw        *string
	normalized *string
	columns    *cronColumns
	next       *nextDateTime
}

type cronColumns struct {
	minutes     *string
	hours       *string
	daysOfMonth *string
	months      *string
	daysOfWeek  *string
}

// // String will return the normalized crontab as a string. Returns string.
// func (cron *Cron) String() string {
// 	return *cron.normalized
// }

// // RawString will return the raw crontab as a string. Returns string.
// func (cron *Cron) RawString() string {
// 	return *cron.raw
// }

// // NextRun will return the Date and time for the next scheduled run. Returns string.
// func (cron *Cron) NextRun() string {
// 	return cron.nextRun()
// }
