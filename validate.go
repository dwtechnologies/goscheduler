package goscheduler

import (
	"fmt"
	"strings"
)

const (
	valueSeparator = ","
	rangeSeparator = "-"
)

var (
	allowedMinutes     = validate{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "26", "27", "28", "29", "30", "31", "32", "33", "34", "35", "36", "37", "38", "39", "40", "41", "42", "43", "44", "45", "46", "47", "48", "49", "50", "51", "52", "53", "54", "55", "56", "57", "58", "59", "*"}
	allowedHours       = validate{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "*"}
	allowedDaysOfMonth = validate{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "26", "27", "28", "29", "30", "31", "*"}
	allowedMonths      = validate{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "*"}
	allowedDaysOfWeek  = validate{"0", "1", "2", "3", "4", "5", "6", "*"}
)

type validate []string

// validateAndSimplify will validate the crontab to make sure that the values are withing the allowed ranges.
// It will also convert ranges to invidual values and set them to the cronColumns struct of the cron struct.
// Returns error.
func (cron *cron) validateAndSimplify() error {
	// Validate minutes.
	err := cron.validateAndSimplifyMinutes()
	if err != nil {
		return err
	}

	// Validate hours.
	err = cron.validateAndSimplifyHours()
	if err != nil {
		return err
	}

	// Validate daysOfMonth.
	err = cron.validateAndSimplifyDaysOfMonth()
	if err != nil {
		return err
	}

	// Validate months.
	err = cron.validateAndSimplifyMonths()
	if err != nil {
		return err
	}

	// Validate daysOfWeek.
	err = cron.validateAndSimplifyDaysOfWeek()
	if err != nil {
		return err
	}

	return nil
}

func (cron *cron) validateAndSimplifyMinutes() error {
	return allowedMinutes.validate(&cron.columns.minutes)
}

func (cron *cron) validateAndSimplifyHours() error {
	return allowedHours.validate(&cron.columns.hours)
}

func (cron *cron) validateAndSimplifyDaysOfMonth() error {
	return allowedDaysOfMonth.validate(&cron.columns.daysOfMonth)
}

func (cron *cron) validateAndSimplifyMonths() error {
	return allowedMonths.validate(&cron.columns.months)
}

func (cron *cron) validateAndSimplifyDaysOfWeek() error {
	return allowedDaysOfWeek.validate(&cron.columns.daysOfWeek)
}

// validate will use the specified validate slice and check that the supplied string is valid. Returns error.
func (validate *validate) validate(raw *string) error {
	slice := strings.Split(*raw, valueSeparator)
	values := []string{}

	// Loop through all the column ranges we have and add them to the values slice.
	for _, value := range slice {
		split := strings.Split(value, rangeSeparator)

		// Validate that the values in split are within the allowed values.
		err := validate.checkIfValid(&split)
		if err != nil {
			return err
		}

		switch {
		// If length of rangeSplit is bigger than 2 there is a syntax error.
		case len(split) > 2:
			return fmt.Errorf("There was an error parsing range %v", value)

		// If length is 2 we have a range. And validate and normalize it.
		case len(split) == 2:
			str, err := validate.createRange(&split[0], &split[1])
			if err != nil {
				return err
			}

			values = append(values, *str)

		// If we didn't have a range, just add the single value.
		default:
			values = append(values, split[0])
		}
	}

	// Set the data of the current column via pointer.
	*raw = strings.Join(values, valueSeparator)
	return nil
}

// checkIfValid will loop through the set values and see if they are valid. Returns error.
func (validate *validate) checkIfValid(slice *[]string) error {
	for _, val := range *slice {
		for _, allowedVal := range *validate {
			if val == allowedVal {
				return nil
			}
		}
	}

	return fmt.Errorf("The cron syntax contains illegal characters")
}
