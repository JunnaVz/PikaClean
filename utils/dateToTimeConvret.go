package utils

import (
	"log"
	"time"
)

func ConvertStringToTime(dateStr string) time.Time {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Fatal(err)
	}
	return t
}
