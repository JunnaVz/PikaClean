// Package cmdUtils provides utility functions for command-line input processing
// in the PikaClean application. It handles continuous input prompts, validation,
// and type conversion for various data types.
package cmdUtils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// InvalidInput is the error message displayed when user input doesn't meet requirements.
const InvalidInput = "ошибка ввода, попробуйте еще раз"

// EndlessReadWord prompts the user for a single word input and continues
// prompting until valid input is received.
//
// Parameters:
//   - requestString: The prompt message to display to the user
//
// Returns:
//   - string: The validated input word with whitespace trimmed
func EndlessReadWord(requestString string) string {
	var input string
	var err error

	fmt.Printf(requestString + ": ")
	for {
		input, err = StringReader(true)
		if err == nil && len(input) > 0 {
			break
		}
		fmt.Print(InvalidInput + ": ")
	}

	input = strings.TrimSpace(input)
	return input
}

// EndlessReadFloat64 prompts the user for a floating-point number and continues
// prompting until a valid float64 is received.
//
// Parameters:
//   - requestString: The prompt message to display to the user
//
// Returns:
//   - float64: The validated floating-point number
func EndlessReadFloat64(requestString string) float64 {
	var input string
	var err error
	var result float64

	fmt.Printf("%s: ", requestString)
	for {
		input, err = StringReader(true)
		if err == nil && len(input) > 0 {
			break
		}
		fmt.Print(InvalidInput + ": ")
	}

	_, err = fmt.Sscanf(input, "%f", &result)
	if err != nil {
		fmt.Print(InvalidInput + ": ")
		return EndlessReadFloat64(requestString)
	}

	return result
}

// EndlessReadInt prompts the user for an integer number and continues
// prompting until a valid int is received.
//
// Parameters:
//   - requestString: The prompt message to display to the user
//
// Returns:
//   - int: The validated integer number
func EndlessReadInt(requestString string) int {
	var input string
	var err error
	var result int

	fmt.Printf("%s: ", requestString)
	for {
		input, err = StringReader(true)
		if err == nil && len(input) > 0 {
			break
		}
		fmt.Print(InvalidInput + ": ")
	}

	_, err = fmt.Sscanf(input, "%d", &result)
	if err != nil {
		fmt.Print(InvalidInput + ": ")
		return EndlessReadInt(requestString)
	}

	return result
}

// EndlessReadRow prompts the user for a full line of text and continues
// prompting until valid input is received.
//
// Parameters:
//   - requestString: The prompt message to display to the user
//
// Returns:
//   - string: The validated input line with whitespace trimmed
func EndlessReadRow(requestString string) string {
	var input string
	var err error

	fmt.Printf("%s: ", requestString)
	for {
		input, err = StringReader(false)
		if err == nil && len(input) > 0 {
			break
		} else {
			fmt.Print(InvalidInput + ": ")
		}
	}

	input = strings.TrimSpace(input)
	return input
}

// stdinReader creates a new buffered reader for standard input.
//
// Returns:
//   - bufio.Reader: A buffered reader for stdin
func stdinReader() bufio.Reader {
	return *bufio.NewReader(os.Stdin)
}

// StringReader reads a string from standard input with options for processing.
//
// Parameters:
//   - firstWordOnly: If true, only the first word of input is returned
//
// Returns:
//   - string: The input string, processed according to the firstWordOnly flag
//   - error: Any error encountered during reading
func StringReader(firstWordOnly bool) (string, error) {
	reader := stdinReader()
	input, err := reader.ReadString('\n')
	input = strings.ReplaceAll(input, "\n", "")

	if firstWordOnly && err == nil {
		words := strings.Fields(input)
		if len(words) > 0 {
			input = words[0]
		} else {
			input = ""
		}
	}

	return input, err
}
