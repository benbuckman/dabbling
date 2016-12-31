package main

import "fmt"

// function that returns a function that returns an int.
func fibonacci() func() int {
	lastTwo := [2]int{}
	i := 0
	return func() int {
		var n int
		if i <= 1 {
			n = i
			lastTwo[i] = n
		} else {
			n = lastTwo[0] + lastTwo[1]
			lastTwo = [2]int{lastTwo[1], n}
		}
		i++
		return n
	}
}

func main() {
	f := fibonacci()
	for i := 0; i < 20; i++ {
		fmt.Println(f())
	}
}
