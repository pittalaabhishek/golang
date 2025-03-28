package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// Prompt the user to enter the file path
	fmt.Print("Enter the file path: ")
	var filePath string
	fmt.Scanln(&filePath)

	// Try to open the file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close() // Ensure the file is closed when the function ends

	// Read and print the contents of the file
	fmt.Println("File contents:")
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	// Check for any errors encountered during reading
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}
