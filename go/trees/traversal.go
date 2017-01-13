package main

import (
	//"fmt"
	. "golang.org/x/tour/tree"
)

// Implement the 4 binary tree traversals described here:
// https://www.cs.cmu.edu/~adamchik/15-121/lectures/Trees/trees.html
// (PreOrder, InOrder, PostOrder, LevelOrder)

func TraversePreOrder(tree *Tree) ([]int) {
	var nums []int
	nums = append(nums, 1)
	return nums
}

func TraversePostOrder(tree *Tree) ([]int) {
	return []int{}
}

func TraverseInOrder(tree *Tree) ([]int) {
	return []int{}
}

func TraverseLevelOrder(tree *Tree) ([]int) {
	return []int{}
}
