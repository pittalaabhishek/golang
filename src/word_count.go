package wordcount

import (
	"regexp"
	"strings"
)

type Frequency map[string]int

func WordCount(input string) Frequency {
	counts := make(Frequency)
	wordPattern := regexp.MustCompile(`\b[\w']+\b`)
	words := wordPattern.FindAllString(input, -1)

	for _, word := range words {
		lowerWord := strings.ToLower(word)
		counts[lowerWord]++
	}

	return counts
}
