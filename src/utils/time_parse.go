package utils

import "time"

func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

func ParseTime(timeStr string) (time.Time, error) {
	return time.Parse("15:04", timeStr)
}

func ParseDateTime(dateStr, timeStr string) (time.Time, error) {
	layout := "2006-01-02 15:04"
	return time.Parse(layout, dateStr+" "+timeStr)
}