package main

import (
	"fmt"
	. "golang.org/x/tour/tree"
)

func incr(n *int) (int) {
	*n++
	return *n
}

func main() {
	n := 0

	t1 := Tree{nil, incr(&n), nil}
	t1.Right = &Tree{nil, incr(&n), nil}
	t1.Right.Right = &Tree{nil, incr(&n), nil}
	t1.Right.Right.Right = &Tree{nil, incr(&n), nil}
	t1.Right.Right.Left = &Tree{nil, incr(&n), nil}

	t1.Left = &Tree{nil, -1, nil}

	fmt.Println(t1)

	// TODO understand: `Tree.New()` uses `rand.Perm()` which doesn't seem to be random at all,
	// the trees are always 1..N :-/

	t2 := New(1)
	fmt.Println(t2)
}
