package main

import "golang.org/x/tour/pic"
//import "fmt"

func Pic(dx, dy int) [][]uint8 {
	//fmt.Printf("dx: %d, dy: %d\n", dx, dy)

	matrix := make([][]uint8, dy)
	for x := range matrix {
		matrix[x] = make([]uint8, dx)

		for y := range matrix[x] {
			matrix[x][y] = uint8((x+y)/2)
		}
	}
	//fmt.Printf("%v", matrix)
	return matrix
}

func main() {
	pic.Show(Pic)
}
//
