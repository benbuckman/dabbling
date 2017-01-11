package main

import "fmt"
import "golang.org/x/tour/tree"

func incr(n *int) (int) {
	*n++
	return *n
}

func main() {
	n := 0

	t1 := tree.Tree{nil, incr(&n), nil}
	t1.Right = &tree.Tree{nil, incr(&n), nil}
	t1.Right.Right = &tree.Tree{nil, incr(&n), nil}
	t1.Right.Right.Right = &tree.Tree{nil, incr(&n), nil}
	t1.Right.Right.Left = &tree.Tree{nil, incr(&n), nil}

	t1.Left = &tree.Tree{nil, -1, nil}

	fmt.Println(t1)

	t2 := tree.New(1)
	fmt.Println(t2)
}
