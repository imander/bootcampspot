package util

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ReadInput(prompt string, s *string) {
	if prompt != "" {
		fmt.Printf("%s [%s]: ", prompt, *s)
	}

	str := getInput()
	if str != "" {
		*s = str
	}
}

func ReadInt(prompt string) int {
	if prompt != "" {
		fmt.Printf("%s: ", prompt)
	}

	str := getInput()
	i, err := strconv.Atoi(str)
	if err != nil {
		fmt.Printf("invalid input: %s", err.Error())
		os.Exit(1)
	}

	return i
}

func getInput() string {
	in := bufio.NewReader(os.Stdin)
	line, err := in.ReadString('\n')
	if err != nil {
		fmt.Printf("error reading std-in: %s", err.Error())
		os.Exit(1)
	}

	return strings.Trim(line, "\r\n")
}
