package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
)

// recoverFromPanic handles unexpected panics gracefully
func recoverFromPanic() {
	if r := recover(); r != nil {
		fmt.Println("Recovered from panic:", r)
	}
}

func main() {
	defer recoverFromPanic() // Ensure recovery from unexpected errors

	// Define a flag that accepts a comma-separated list of numbers
	input := flag.String("numbers", "", "Comma-separated list of numbers")
	flag.Parse()

	if *input == "" {
		fmt.Println("Error: No numbers provided. Use -numbers flag.")
		os.Exit(1)
	}

	nums := parseNumbers(*input)
	sum := 0
	for _, num := range nums {
		sum += num
	}

	fmt.Println("Sum of valid numbers:", sum)
}

func parseNumbers(input string) []int {
	var numbers []int
	values := splitString(input)

	for _, val := range values {
		num, err := convertToInt(val)
		if err != nil {
			fmt.Printf("Skipping invalid input '%s': %v\n", val, err)
			continue
		}
		numbers = append(numbers, num)
	}
	return numbers
}

func splitString(input string) []string {
	return append([]string{input}, flag.Args()...)
}

func convertToInt(value string) (int, error) {
	num, err := strconv.Atoi(value)
	if err != nil {
		return 0, errors.New("not a valid integer")
	}
	return num, nil
}
