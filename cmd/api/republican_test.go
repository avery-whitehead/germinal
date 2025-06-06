package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/rickb777/date"
)

func TestGetDayMonthYear(t *testing.T) {
	var tests = []struct {
		date   date.Date
		expect *RepublicanDate
	}{
		{date.New(2025, time.January, 1),
			&RepublicanDate{
				Year:     233,
				DayOrd:   12,
				MonthOrd: 4,
			}},
		{date.New(2024, time.May, 25),
			&RepublicanDate{
				Year:     232,
				DayOrd:   7,
				MonthOrd: 9,
			}},
		{date.New(1792, time.September, 22),
			&RepublicanDate{
				Year:     1,
				DayOrd:   1,
				MonthOrd: 1,
			}},
		{date.New(2024, time.September, 18),
			&RepublicanDate{
				Year:     232,
				DayOrd:   3,
				MonthOrd: 13,
			}},
		{date.New(2023, time.September, 18),
			&RepublicanDate{
				Year:     231,
				DayOrd:   2,
				MonthOrd: 13,
			}},
		{date.New(9999, time.December, 31),
			&RepublicanDate{
				Year:     8208,
				DayOrd:   11,
				MonthOrd: 4,
			}},
	}

	for _, tt := range tests {
		testname := tt.date.UTC().String()
		t.Run(testname, func(t *testing.T) {
			ans := getDayMonthYear(tt.date)
			if ans.DayOrd != tt.expect.DayOrd || ans.MonthOrd != tt.expect.MonthOrd || ans.Year != tt.expect.Year {
				t.Errorf("got %+v, expected %+v", ans, tt.expect)
			}
		})
	}
}

func TestIsLeapyear(t *testing.T) {
	var tests = []struct {
		year   int
		expect bool
	}{
		{1, false},
		{3, true},
		{4, false},
		{16, false},
		{17, false},
		{20, true},
		{230, false},
		{232, true},
		{260, true},
		{300, false},
		{400, true},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%d", tt.year)
		t.Run(testname, func(t *testing.T) {
			ans := isLeapYear(tt.year)
			if ans != tt.expect {
				t.Errorf("got %t, expected %t", ans, tt.expect)
			}
		})
	}
}

func TestYearToRoman(t *testing.T) {
	var tests = []struct {
		year   int
		expect string
	}{
		{2025, "MMXXV"},
		{1, "I"},
		{3999, "MMMCMXCIX"},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%d", tt.year)
		t.Run(testname, func(t *testing.T) {
			ans := yearToRoman(tt.year)
			if ans != tt.expect {
				t.Errorf("got %s, expected %s", ans, tt.expect)
			}
		})
	}
}
