package utils

import "time"

func FormatDate(t time.Time) string {
	return t.Format("02-01-2006") // DD-MM-YYYY
}
