package timeparser

import (
	"fmt"
	"time"
)

func ParseMonthYear(monthYear string) (time.Time, error) {
	parsed, err := time.Parse("01-2006", monthYear)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format, expected MM-YYYY")
	}
	return time.Date(parsed.Year(), parsed.Month(), 1, 0, 0, 0, 0, time.UTC), nil
}
