package goscheduler

import (
	"fmt"
	"strings"
)

const (
	columnSeparator = " "
)

var (
	conversionMonths     = conversion{"JAN": "1", "FEB": "2", "MAR": "3", "APR": "4", "MAY": "5", "JUN": "6", "JUL": "7", "AUG": "8", "SEP": "9", "OCT": "10", "NOV": "11", "DEC": "12"}
	conversionDaysOfWeek = conversion{"SUN": "0", "MON": "1", "TUE": "2", "WED": "3", "THU": "4", "FRI": "5", "SAT": "6", "7": "0"}
	conversionMacros     = conversion{"@yearly": "0 0 1 1 *", "@annually": "0 0 1 1 *", "@monthly": "0 0 1 * *", "@weekly": "0 0 * * 0", "@daily": "0 0 * * *", "@hourly": "0 * * * *", "@minutely": "* * * * *", "@secondly": "* * * * * *"}
)

type conversion map[string]string

// normalizeCron will generate a standard cron syntax according to POSIX (+ optional extra column in the begginning for seconds).
// Will return a normalized cron string and error.
func (cron *cron) normalizeCron() (*string, error) {
	// Replace any macros with a real cron data string. If it did replace we don't need to do anymore parsing. Just returned the changed cron-string.
	normalized := cron.convertMarcos()
	if *cron.raw != *normalized {
		return normalized, nil
	}

	slice := strings.Split(*cron.raw, columnSeparator)
	if len(slice) != 5 {
		return nil, fmt.Errorf("Cron syntax error. Current syntax is \"%v\" and has %v columns. Must be 5 columns seperated by <space>", *cron.raw, len(slice))
	}
	cron.columns = &cronColumns{minutes: &slice[0], hours: &slice[1], daysOfMonth: &slice[2], months: cron.convertMonths(&slice[3]), daysOfWeek: cron.convertDaysOfWeek(&slice[4])}

	// Validate the cron values.
	err := cron.validateAndSimplify()
	if err != nil {
		return nil, err
	}

	// Create cron from the converted and validated slice.
	cronstring := fmt.Sprintf("%v %v %v %v %v", *cron.columns.minutes, *cron.columns.hours, *cron.columns.daysOfMonth, *cron.columns.months, *cron.columns.daysOfWeek)
	return &cronstring, nil
}
