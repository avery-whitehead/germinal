package internal

import (
	"errors"
	"slices"
	"time"
)

var (
	RepublicEraStart = time.Date(1792, time.September, 22, 0, 0, 0, 0, time.UTC)
	RepublicEraEnd   = time.Date(1811, time.September, 23, 0, 0, 0, 0, time.UTC)

	ErrBeforeCalendar = errors.New("date is before year one of the republic")
	ErrDateTooHigh    = errors.New("date is after the year 9999")
)

type RepublicanDate struct {
	Year       int
	YearRoman  string
	Month      string
	MonthOf    string
	MonthOrd   int
	Day        string
	DayOrd     int
	Dedication string
}

// Converts a Gregorian (ISO) date to its Republican equivalent
func toRepublican(date time.Time) (*RepublicanDate, error) {
	if date.Before(RepublicEraStart) {
		return nil, ErrBeforeCalendar
	}

	if date.After(time.Date(9999, time.December, 31, 0, 0, 0, 0, time.UTC)) {
		return nil, ErrDateTooHigh
	}

	delta := (date.Sub(RepublicEraStart).Hours() / 24) + 1

	year, start := 1, 1

	// Count the amount of years (including leap years) that
	// have elapsed in the days between the target date and the start
	// of the Republican
	for {
		var end int
		if isLeapYear(year) {
			end = start + 365
		} else {
			end = start + 364
		}

		if end >= int(delta) {
			break
		}

		year += 1
		start = end + 1
	}

	var repubDate RepublicanDate
	repubDate.Year = year

	dayInYear := int(delta) - start
	// Each month is 30 days
	repubDate.MonthOrd = (dayInYear / 30) + 1
	repubDate.DayOrd = (dayInYear % 30) + 1

	return &repubDate, nil
}

// Returns true if the given Republican year is a leap year
// Leap yers after the first four were left undefined, so this uses the
// reformed calendar proposal to make leap years regular and predictable
func isLeapYear(year int) bool {
	knownLeapYears := []int{3, 7, 11, 15}
	if year <= 16 {
		return slices.Contains(knownLeapYears, year)
	}

	// Leap years are divisible by 4, unless they are also divisible
	// by 100 and not divisible by 400
	if year%4 == 0 {
		if year%100 == 0 && year%400 != 0 {
			return false
		}
		return true
	}
	return false
}
