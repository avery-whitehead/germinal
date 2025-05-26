package internal

import (
	"fmt"
	"testing"
	"time"
)

func TestToRepublican(t *testing.T) {
	var tests = []struct {
		date   time.Time
		expect *RepublicanDate
	}{
		{time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
			&RepublicanDate{
				Year:     233,
				DayOrd:   12,
				MonthOrd: 4,
			}},
		{time.Date(2024, time.May, 25, 0, 0, 0, 0, time.UTC),
			&RepublicanDate{
				Year:     232,
				DayOrd:   7,
				MonthOrd: 9,
			}},
		{time.Date(1792, time.September, 22, 0, 0, 0, 0, time.UTC),
			&RepublicanDate{
				Year:     1,
				DayOrd:   1,
				MonthOrd: 1,
			}},
		{time.Date(2024, time.September, 18, 0, 0, 0, 0, time.UTC),
			&RepublicanDate{
				Year:     232,
				DayOrd:   3,
				MonthOrd: 13,
			}},
		{time.Date(2023, time.September, 18, 0, 0, 0, 0, time.UTC),
			&RepublicanDate{
				Year:     231,
				DayOrd:   2,
				MonthOrd: 13,
			}},
		{time.Date(9999, time.December, 31, 0, 0, 0, 0, time.UTC),
			&RepublicanDate{
				Year:     8208,
				DayOrd:   11,
				MonthOrd: 4,
			}},
	}

	for _, tt := range tests {
		testname := tt.date.UTC().String()
		t.Run(testname, func(t *testing.T) {
			ans, _ := toRepublican(tt.date)
			if ans.DayOrd != tt.expect.DayOrd || ans.MonthOrd != tt.expect.MonthOrd || ans.Year != tt.expect.Year {
				t.Errorf("got %+v, expected %+v", ans, tt.expect)
			}
		})
	}
}

func TestToRepublicanErrors(t *testing.T) {
	var tests = []struct {
		date   time.Time
		expect error
	}{
		{time.Date(1792, time.September, 21, 0, 0, 0, 0, time.UTC), ErrBeforeCalendar},
		{time.Date(10000, time.January, 1, 0, 0, 0, 0, time.UTC), ErrDateTooHigh},
	}

	for _, tt := range tests {
		testname := tt.date.UTC().String()
		t.Run(testname, func(t *testing.T) {
			_, err := toRepublican(tt.date)
			if err != tt.expect {
				t.Errorf("got %s, expected %s", err, tt.expect)
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
