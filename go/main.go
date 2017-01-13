package main

import (
	"fmt"
	"strconv"
)

var topNumber int // shared top-level scope

func multiply(x, y float32) (z float32) {
	z = x * y
	return
}

func setNumber(n int) {
	topNumber = n
}

func incrementPtr(n *int) *int {
	*n = *n + 1
	return n
}

type person struct {
	name string // no way to set default value directly
	age  int
}

// pseudo-constructor ?
func Person(name string, age int) person {
	if len(name) == 0 {
		name = "who?"
	}
	return person{name: name, age: age}
}

func main() {
	fmt.Println(multiply(10, 3.5))

	fmt.Printf("Initial num: %v\n", topNumber)
	defer fmt.Printf("Deferred original num: %v\n", topNumber) // evaluates w/ current value (0)!
	setNumber(25)
	fmt.Printf("Set num: %v\n", topNumber)

	fmt.Printf("Type: %T Value: %v\n", topNumber, topNumber)

	n := topNumber
	numRange1 := make([]int, n)
	var numRange2 []int

	for i := 1; i <= n; i++ {
		numRange1[i-1] = i
		numRange2 = append(numRange2, i)

		var char string
		if isEven := (i % 2) == 0; isEven {
			char = "-"
		} else {
			char = "+"
		}
		fmt.Printf("%v ", i)
		for j := 1; j <= i; j++ {
			fmt.Printf(char)
		}
		fmt.Printf("\n")
	}

	fmt.Println(numRange1)
	fmt.Println(numRange2)

	n2 := incrementPtr(&topNumber)
	fmt.Printf("Incremented: %v\n", topNumber)
	fmt.Printf("Copied pointer?: %v\n", *n2)

	//// Doesn't work -- TODO understand
	//incrementPtr(&n2)
	//fmt.Printf("increment again: %v %v", *topNumber, *n2)

	me := Person("Ben", 31)
	fmt.Printf("Me: %T %v\n", me, me)
	fmt.Println("Hello", me.name)

	someone := Person("", 0)
	fmt.Printf("Someone: %T %v\n", someone, someone)

	numberMap1 := make(map[string]int)
	numberMap1["thing 1"] = 100
	numberMap1["thing 2"] = 200
	fmt.Println(numberMap1)

	numberMap2 := map[string]int{
		"thing 3": 300,
		"thing 4": 400,
		"thing 5": 500,
	}
	delete(numberMap2, "thing 4")
	fmt.Println(numberMap2)

	for i := 1; i <= 10; i++ {
		key := "thing " + strconv.Itoa(i)
		_, exists := numberMap2[key]
		fmt.Printf("'%s' exists in numberMap2? %v\n", key, exists)
	}

	fmt.Println(len("hello"))

	fmt.Println("---------------------------------")

	var b byte = 100
	fmt.Println(b)
	s := strconv.Itoa(int(b))
	fmt.Println(s)

	bytes := []byte{1, 2, 3, 4}
	fmt.Println(bytes)
	var strings []string
	for _, b := range bytes {
		strings = append(strings, strconv.Itoa(int(b)))
	}
	fmt.Println(strings)

	fmt.Println("---------------------------------")

	// "shift" an array
	nums := []int{}
	for i := 0; i < 10; i++ {
		nums = append(nums, i)
	}
	for len(nums) > 0 {
		nums = nums[1:]
		fmt.Println(nums)
	}

}
