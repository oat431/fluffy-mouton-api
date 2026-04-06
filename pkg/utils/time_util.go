package utils

import "time"

func GetCurrentTime() string {
	return time.Now().Format(time.RFC3339)
}

func GetTimeFromString(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}
