package goscheduler

// setField will set the value of the corresponding field in nextDateTime struct.
func (cron *cron) setField(t *int, s *string, field *string, extra *bool, future *bool) {
	switch *field {
	case "minute":
		cron.setFieldMinute(t, s, field, extra, future)

	case "hour":
		cron.setFieldHour(t, s, field, extra, future)

	case "dayOfMonth":
		cron.setFieldDayOfMonth(t, s, field, extra, future)

	case "month":
		cron.setFieldMonth(t, s, field, extra, future)
	}
}

func (cron *cron) setFieldMinute(t *int, s *string, field *string, extra *bool, future *bool) {
	if cron.next.futureHour || cron.next.futureDayOfMonth || cron.next.futureMonth {
		cron.next.minute = *t
		return
	}
	if cron.next.futureRun {
		return
	}

	cron.next.minute += *t
	if *extra {
		cron.next.hour++
	}
}

func (cron *cron) setFieldHour(t *int, s *string, field *string, extra *bool, future *bool) {
	if cron.next.futureHour || cron.next.futureDayOfMonth || cron.next.futureMonth {
		cron.next.hour = *t
		return
	}
	if cron.next.futureRun {
		return
	}

	cron.next.hour += *t
	if *extra {
		cron.next.minute = 0
		cron.next.dayOfMonth++
	}
	if *future {
		cron.next.futureHour = true
		cron.next.future = true
	}
}

func (cron *cron) setFieldDayOfMonth(t *int, s *string, field *string, extra *bool, future *bool) {
	if cron.next.futureDayOfMonth || cron.next.futureMonth {
		cron.next.dayOfMonth = *t
		return
	}
	if cron.next.futureRun {
		return
	}

	cron.next.dayOfMonth += *t
	if *extra {
		cron.next.minute = 0
		cron.next.hour = 0
		cron.next.month++
	}
	if *future {
		cron.next.futureDayOfMonth = true
		cron.next.future = true
	}
}

func (cron *cron) setFieldMonth(t *int, s *string, field *string, extra *bool, future *bool) {
	if cron.next.futureMonth {
		cron.next.month = *t
		return
	}
	if cron.next.futureRun {
		return
	}

	cron.next.month += *t
	if *extra {
		cron.next.minute = 0
		cron.next.hour = 0
		cron.next.dayOfMonth = 1
		cron.next.year++
	}
	if *future {
		cron.next.futureMonth = true
		cron.next.future = true
	}
}
