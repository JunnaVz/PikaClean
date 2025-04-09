package cmdUtils

import "unicode/utf8"

func TruncateString(str string, num int) string {
	if utf8.RuneCountInString(str) <= num {
		return str
	}
	i := 0
	for j := range str {
		if i == num {
			return str[:j] + "..."
		}
		i++
	}
	return str
}
