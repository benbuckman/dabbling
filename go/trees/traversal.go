package main

import (
	//"fmt"
	. "golang.org/x/tour/tree"
)

// Implement the 4 binary tree traversals described here:
// https://www.cs.cmu.edu/~adamchik/15-121/lectures/Trees/trees.html
// (PreOrder, InOrder, PostOrder, LevelOrder)

// PreOrder traversal - visit the parent first and then left and right children
func TraversePreOrder(tree *Tree) ([]int) {
	var values []int
	var recursivelyTraverse func(node *Tree)
	recursivelyTraverse = func(node *Tree) {
		values = append(values, node.Value)
		if node.Left != nil {
			recursivelyTraverse(node.Left)
		}
		if node.Right != nil {
			recursivelyTraverse(node.Right)
		}
	}
	recursivelyTraverse(tree)
	return values
}

// PostOrder traversal - visit left child, then the right child and then the parent
func TraversePostOrder(tree *Tree) ([]int) {
	var values []int
	var recursivelyTraverse func(node *Tree)
	recursivelyTraverse = func(node *Tree) {
		if node.Left != nil {
			recursivelyTraverse(node.Left)
		}
		if node.Right != nil {
			recursivelyTraverse(node.Right)
		}
		values = append(values, node.Value)
	}
	recursivelyTraverse(tree)
	return values
}

// InOrder traversal - visit the left child, then the parent and the right child
func TraverseInOrder(tree *Tree) ([]int) {
	var values []int
	var recursivelyTraverse func(node *Tree)
	recursivelyTraverse = func(node *Tree) {
		if node.Left != nil {
			recursivelyTraverse(node.Left)
		}
		values = append(values, node.Value)
		if node.Right != nil {
			recursivelyTraverse(node.Right)
		}
	}
	recursivelyTraverse(tree)
	return values
}

// LevelOrder traversal (breadth-first) - visits nodes by levels from top to bottom and from left to right
func TraverseLevelOrder(tree *Tree) ([]int) {
	var values []int
	q := []*Tree{tree}

	for len(q) > 0 {
		node := q[0]
		q = q[1:]
		values = append(values, node.Value)
		if node.Left != nil {
			q = append(q, node.Left)
		}
		if node.Right != nil {
			q = append(q, node.Right)
		}
	}

	return values
}
