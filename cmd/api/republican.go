package main

import (
	"errors"
	"slices"
	"time"

	"github.com/rickb777/date"
)

var (
	RepublicEraStart = date.New(1792, time.September, 22)
	MaxAllowedDate   = date.New(9999, time.December, 31)

	ErrBeforeCalendar = errors.New("date is before year one of the republic")
	ErrDateTooHigh    = errors.New("date is after the year 9999")

	KnownLeapYears = []int{3, 7, 11, 15}
	RomanNumerals  = []roman{
		{1000, "M"},
		{900, "CM"},
		{500, "D"},
		{400, "CD"},
		{100, "C"},
		{90, "XC"},
		{50, "L"},
		{40, "XL"},
		{10, "X"},
		{9, "IX"},
		{5, "V"},
		{4, "IV"},
		{1, "I"},
	}
)

type RepublicanDate struct {
	Year       int    `json:"year"`
	YearRoman  string `json:"yearRoman"`
	Month      string `json:"month"`
	MonthOf    string `json:"monthOf"`
	MonthOrd   int    `json:"monthOrd"`
	Day        string `json:"day"`
	DayOrd     int    `json:"dayOrd"`
	Dedication string `json:"dedication"`
}

type DayMonthYear struct {
	DayOrd   int
	MonthOrd int
	Year     int
}

type roman struct {
	number int
	value  string
}

// Converts a Gregorian (ISO) date to its Republican equivalent
func (app *application) toRepublican(date date.Date) (*RepublicanDate, error) {
	if date.Before(RepublicEraStart) {
		return nil, ErrBeforeCalendar
	}

	if date.After(MaxAllowedDate) {
		return nil, ErrDateTooHigh
	}

	dateValues := getDayMonthYear(date)

	// Get string values from database
	dbValues, err := app.db.GetDbValues(dateValues.DayOrd, dateValues.MonthOrd)
	if err != nil {
		return nil, err
	}

	return &RepublicanDate{
		Year:       dateValues.Year,
		YearRoman:  yearToRoman(dateValues.Year),
		Month:      dbValues.Month,
		MonthOf:    dbValues.MonthOf,
		MonthOrd:   dateValues.MonthOrd,
		Day:        dbValues.Day,
		DayOrd:     dateValues.DayOrd,
		Dedication: dbValues.Dedication,
	}, nil
}

func getDayMonthYear(date date.Date) *DayMonthYear {
	delta := date.Sub(RepublicEraStart)

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

		if end >= int(delta+1) {
			break
		}

		year += 1
		start = end + 1
	}

	dayInYear := int(delta+1) - start

	return &DayMonthYear{
		// Each month is 30 days
		MonthOrd: (dayInYear / 30) + 1,
		DayOrd:   (dayInYear % 30) + 1,
		Year:     year,
	}
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

func yearToRoman(year int) string {
	if year >= 4000 {
		return "-"
	}
	var result string

	for _, v := range RomanNumerals {
		for year >= v.number {
			result += v.value
			year = year - v.number
		}
	}

	return result
}
