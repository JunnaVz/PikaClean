package cmdUtils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const InvalidInput = "ошибка ввода, попробуйте еще раз"

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

func stdinReader() bufio.Reader {
	return *bufio.NewReader(os.Stdin)
}

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
