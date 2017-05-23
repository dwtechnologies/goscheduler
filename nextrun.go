package goscheduler

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	wildCard = "*"
)

type nextDateTime struct {
	minute     int
	hour       int
	dayOfMonth int
	month      int
	dayOfWeek  int
	year       int
	time       *time.Time

	future           bool
	futureRun        bool
	futureHour       bool
	futureDayOfMonth bool
	futureMonth      bool
}

func (job *Job) genNextRun() *int64 {
	now := time.Now()
	genTime := now.Add(time.Duration(60-now.Second()) * time.Second)
	job.cron.next = &nextDateTime{
		time: &genTime,
		year: genTime.Year(),
	}
	job.next()

	nextRun, _ := time.ParseInLocation(format, fmt.Sprintf("%v-%v-%v %v:%v", job.cron.next.year, *job.cron.convertAndAppendMonth(), *job.cron.convertAndAppendDayOfMonth(), *job.cron.convertAndAppendHour(), *job.cron.convertAndAppendMinute()), time.Local)
	diff := nextRun.Unix() - now.Unix()

	return &diff
}

func (cron *cron) convertAndAppendMinute() *string {
	return cron.convertAndAppend(strconv.Itoa(cron.next.minute))
}

func (cron *cron) convertAndAppendHour() *string {
	return cron.convertAndAppend(strconv.Itoa(cron.next.hour))
}

func (cron *cron) convertAndAppendDayOfMonth() *string {
	return cron.convertAndAppend(strconv.Itoa(cron.next.dayOfMonth))
}

func (cron *cron) convertAndAppendMonth() *string {
	return cron.convertAndAppend(strconv.Itoa(cron.next.month))
}

func (*cron) convertAndAppend(str string) *string {
	if len(str) == 1 {
		str = fmt.Sprintf("0%v", str)
		return &str
	}
	return &str
}

func (job *Job) next() {
	fields := []string{"minute", "hour", "month", "dayOfMonth"}
	for _, field := range fields {
		job.cron.nextField(&field)
	}

	if !job.cron.next.futureRun {
		job.cron.next.futureRun = true
		job.next()
	}
}

func (cron *cron) nextField(field *string) {
	val := ""
	now := 0

	switch *field {
	case "minute":
		val = strings.Split(cron.normalized, columnSeparator)[0]
		now = cron.next.time.Minute()
	case "hour":
		val = strings.Split(cron.normalized, columnSeparator)[1]
		now = cron.next.time.Hour()
	case "month":
		val = strings.Split(cron.normalized, columnSeparator)[3]
		_, m, _ := cron.next.time.Date()
		now = int(m)
	case "dayOfMonth":
		val = strings.Split(cron.normalized, columnSeparator)[2]
		// Add weekdays to the list of dayOfMonth.
		val = *(cron).addWeekdaysToDaysOfMonth(&val)
		now = cron.next.time.Day()
	}

	cron.getNextFromValues(&now, &val, field)
}

func (cron *cron) getNextFromValues(t *int, s *string, field *string) {
	if *s == wildCard {
		cron.wildcard(t, s, field)
		return
	}
	sliceint := cron.createSliceAndSort(s)

	if cron.next.futureRun && cron.next.futureMonth {
		extra := false
		future := false
		cron.setField(&(*sliceint)[0], s, field, &extra, &future)
		return
	}

	for _, val := range *sliceint {
		if val >= *t {
			extra := false
			future := true
			cron.setField(&val, s, field, &extra, &future)
			return
		}
	}

	extra := true
	future := true
	cron.setField(&(*sliceint)[0], s, field, &extra, &future)
}

func (cron *cron) wildcard(t *int, s *string, field *string) {
	val := *t
	extra := false
	future := false

	if cron.next.futureRun {
		if *field == "month" || *field == "dayOfMonth" {
			val = 1
		}
	}

	cron.setField(&val, s, field, &extra, &future)
}

func (*cron) createSliceAndSort(s *string) *[]int {
	slice := strings.Split(*s, valueSeparator)
	sliceint := []int{}

	for _, val := range slice {
		conv, _ := strconv.Atoi(val)
		sliceint = append(sliceint, conv)
	}
	sort.Ints(sliceint)

	return &sliceint
}
