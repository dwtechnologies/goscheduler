package goscheduler

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	days      = 7
	secsInDay = 86400
)

// addWeekdaysToDaysOfMonth will convert any weekdays to daysOfMonth format and add it to daysOfMonth list.
// Returns a combined string with daysOfMonth and weekdays.
func (cron *cron) addWeekdaysToDaysOfMonth(org *string) *string {
	s := strings.Split(cron.normalized, columnSeparator)[4]
	switch {
	case *org == wildCard && s == wildCard:
		w := wildCard
		return &w

	case s == wildCard:
		return org
	}

	slice := []string{}

	for _, val := range *cron.createSliceAndSort(&s) {
		slice = append(slice, *(cron).addWeekdays(&val))
	}

	if *org == wildCard {
		str := strings.Join(slice, ",")
		return &str
	}
	str := fmt.Sprintf("%v,%v", *org, strings.Join(slice, ","))
	return &str
}

// addWeekdays will calculate which dayOfMonth the specified weekday is.
func (cron *cron) addWeekdays(to *int) *string {
	date := *cron.next.time
	today := int(cron.next.time.Weekday())
	diff := *to - today

	// If we started from a new month set current day to 01.
	// Since we need to start parsing from day 01 in these cases.
	if cron.next.future {
		if cron.next.futureMonth {
			date, _ = time.Parse("2006-01-02", fmt.Sprintf("%v-%v-01", cron.next.year, cron.convertAndAppendMonth()))
		}
		today = int(date.Weekday())
		diff = *to - today
	}

	if diff < 0 {
		diff = days + diff
	}

	str := strconv.Itoa(date.Add(time.Duration(secsInDay*(diff)) * time.Second).Day())
	return &str
}
