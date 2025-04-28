// Package cmdUtils provides utility functions for command-line input processing
// and string manipulation for the PikaClean application. It includes functions
// for reading user input with validation and formatting text for display.
package cmdUtils

import "unicode/utf8"

// TruncateString shortens a string to a specified maximum length,
// adding an ellipsis ("...") if the string exceeds that length.
// The function properly handles UTF-8 encoded strings by counting
// characters rather than bytes.
//
// Parameters:
//   - str: The input string to be truncated
//   - num: Maximum number of characters to keep before truncation
//
// Returns:
//   - string: The truncated string with ellipsis if needed, or the original string
//     if it was already shorter than the specified length
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
