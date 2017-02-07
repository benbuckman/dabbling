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
	nums := make([]int, 100)
	for i := 0; i < 100; i++ {
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
	for len(nums) > 0 {
		iterations++
		ind = int(math.Floor((float64(len(nums)) - 1) / 2))	// midpoint
		if nums[ind] == n {
			return
		} else if nums[ind] > n {
			nums = nums[0:ind]
		} else {
			nums = nums[ind+1:]
		}
		fmt.Println(nums)
	}
	err = fmt.Errorf("Could not find %v", n)
	return
}

func main() {
	nums := randomNumbers()
	fmt.Println(nums)
	numToFind := randomNumberInList(nums)
	fmt.Println("finding", numToFind)
	ind, iterations, err := binarySearch(nums, numToFind)
	if err != nil {
		fmt.Println("Failed:", err)
	} else {
		fmt.Printf("Found %v at index %v in %v iterations\n", numToFind, ind, iterations)
	}
}
