package goscheduler

import "strings"

// convertMarcos will convert common non-standard cron macros into a cron line. Returns string.
func (cron *cron) convertMarcos() *string {
	return conversionMacros.conversion(&cron.raw)
}

// convertDaysOfWeek will convert character based version of days of the week into numeric. Returns string.
func (*cron) convertDaysOfWeek(days *string) *string {
	return conversionDaysOfWeek.conversion(days)
}

// convertMonths will convert character based version of months into numeric. Returns string.
func (*cron) convertMonths(months *string) *string {
	return conversionMonths.conversion(months)
}

// conversion will use the specified conversion struct and convert any found strings.
func (conv *conversion) conversion(raw *string) *string {
	data := *raw

	for key, value := range *conv {
		data = strings.Replace(data, key, value, -1)
	}

	return &data
}
