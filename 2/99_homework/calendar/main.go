package main

import "time"

type Calendar struct {
	year  int
	month time.Month
	day   int
}

func NewCalendar(t time.Time) Calendar {
	return Calendar{t.Year(), t.Month(), t.Day()}
}

func (cal Calendar) CurrentQuarter() int {
	switch cal.month {
	case 1, 2, 3:
		return 1
	case 4, 5, 6:
		return 2
	case 7, 8, 9:
		return 3
	case 10, 11, 12:
		return 4
	}
	return 1
}
