package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
	"math"
)

// Implement a simple binary search algorithm.

// build a random list of 100 numbers from 0..999
func randomNumbers() ([]int) {
	var rander *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	l := 1000
	nums := make([]int, l)
	for i := 0; i < l; i++ {
		nums[i] = rander.Intn(1000)
	}
	sort.Ints(nums)
	return nums
}

// pick a random number in the list
func randomNumberInList(nums []int) (int) {
	var rander *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	ind := rander.Intn(len(nums))
	return nums[ind]
}

func binarySearch(nums []int, n int) (ind, iterations int, err error) {
	var bottom, top int = 0, len(nums) - 1
	for top >= bottom {
		iterations++
		ind = bottom + int(math.Floor(float64((top - bottom) / 2)))	// midpoint
		fmt.Printf("%v..(%v)..%v (%v)\n", bottom, ind, top, nums[ind])
		if nums[ind] == n {
			return
		} else if nums[ind] < n {
			bottom = ind + 1
		} else {
			top = ind - 1
		}
	}
	err = fmt.Errorf("Could not find %v", n)
	return
}

func main() {
	nums := randomNumbers()
	fmt.Println(nums)
	numToFind := randomNumberInList(nums)
	fmt.Println("Finding", numToFind)
	ind, iterations, err := binarySearch(nums, numToFind)
	if err != nil {
		fmt.Println("Failed:", err)
	} else {
		fmt.Printf("Found %v at index %v in %v iterations\n", numToFind, ind, iterations)
	}
}
