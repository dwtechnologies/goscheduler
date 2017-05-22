package goscheduler

import (
	"fmt"
	"strconv"
	"strings"
)

// createRange takes range1 string and range2 string and will check that the range is valid and creates it.
// Will return . Returns &string and error.
func (validate *validate) createRange(str1 *string, str2 *string) (*string, error) {
	err := validate.checkRange(str1, str2)
	if err != nil {
		return nil, err
	}

	strRange, err := validate.rangeString(str1, str2)
	if err != nil {
		return nil, err
	}

	return strRange, nil
}

// convertToInt will convert two strings to int. Returns int, int and error.
func (*validate) convertToInt(str1 *string, str2 *string) (*int, *int, error) {
	int1, err := strconv.Atoi(*str1)
	if err != nil {
		return nil, nil, fmt.Errorf("Couldn't convert value %v to integer", *str1)
	}

	int2, err := strconv.Atoi(*str2)
	if err != nil {
		return nil, nil, fmt.Errorf("Couldn't convert value %v to integer", *str2)
	}

	return &int1, &int2, nil
}

// checkRange will make sure that left hand numeric is not larger than the right hand. Returns error.
func (validate *validate) checkRange(range1 *string, range2 *string) error {
	int1, int2, err := validate.convertToInt(range1, range2)
	if err != nil {
		return err
	}

	if *int1 > *int2 {
		return fmt.Errorf("Left hand integer %v is larger than right hand integer %v. Invalid range", *int1, *int2)
	}

	return nil
}

// rangeString will generate the range between int1 and int2. Returns string.
func (validate *validate) rangeString(str1 *string, str2 *string) (*string, error) {
	int1, int2, err := validate.convertToInt(str1, str2)
	if err != nil {
		return nil, err
	}

	slice := []string{}
	for i := *int1; i < *int2+1; i++ {
		slice = append(slice, strconv.Itoa(i))
	}

	str := strings.Join(slice, valueSeparator)
	return &str, nil
}
