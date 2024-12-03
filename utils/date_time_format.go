package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func SplitDateTime(dt time.Time) (string, string) {
	date := dt.Format("2006-01-02")
	ampm := "AM"
	hour := dt.Hour()

	if hour >= 12 {
		ampm = "PM"
		if hour > 12 {
			hour -= 12
		}
	}
	if hour == 0 {
		hour = 12
	}

	time := fmt.Sprintf("%s-%02d-%02d", ampm, hour, dt.Minute())
	return date, time
}

func CombineDateTime(date string, timeStr string) (time.Time, error) {
	parts := strings.Split(timeStr, "-")
	ampm, hourStr, minStr := parts[0], parts[1], parts[2]

	hour, _ := strconv.Atoi(hourStr)
	min, _ := strconv.Atoi(minStr)

	if ampm == "PM" && hour != 12 {
		hour += 12
	}
	if ampm == "AM" && hour == 12 {
		hour = 0
	}

	dateTimeStr := fmt.Sprintf("%s %02d:%02d:00", date, hour, min)
	return time.Parse("2006-01-02 15:04:05", dateTimeStr)
}
