package main

import (
	"fmt"
	"golang.org/x/tour/wc"
	"strings"
)

func WordCount(s string) map[string]int {
	counts := map[string]int{}

	words := strings.Fields(s)
	fmt.Println(words)

	for _, word := range words {
		_, exists := counts[word]
		if !exists {
			counts[word] = 0
		}
		counts[word]++
	}

	return counts
}

func main() {
	wc.Test(WordCount)
}
