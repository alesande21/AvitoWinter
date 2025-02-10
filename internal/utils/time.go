package utils

import "time"

func GetCurrentTimeRFC3339() string {
	currentTime := time.Now()

	formattedTime := currentTime.Format(time.RFC3339)
	return formattedTime
}

func GetCurrentTime() time.Time {
	currentTime := time.Now()
	return currentTime
}

func ToFormatRFC3339(t time.Time) string {
	formattedTime := t.Format(time.RFC3339)
	return formattedTime
}
