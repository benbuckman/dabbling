package main

import "fmt"

type Thing struct {
	Name string
	NameRef *string
}

func parseThing(th *Thing) {
	fmt.Println(th)
	fmt.Println(*th)
	fmt.Println((*th).Name)
	fmt.Println(th.Name)	// syntactic sugar for previous
}

func main() {
	name := "dog"

	th := Thing{
		Name: "monkey",
		NameRef: &name,
	}
	fmt.Println(th)
	fmt.Println(*th.NameRef)

	thPtr := &th
	fmt.Println(thPtr.Name)
	fmt.Println(thPtr.Name)
	fmt.Println(thPtr.NameRef)

	parseThing(&th)
}
