package main

import (
	"testing"
	. "golang.org/x/tour/tree"
	"reflect"
)

func exampleTree()(Tree) {
	return Tree{
		Left: &Tree{
			Left: &Tree{
				Left: nil,
				Value: 9,
				Right: nil,
			},
			Value: 5,
			Right: &Tree{
				Left: &Tree{
					Left: nil,
					Value: 1,
					Right: nil,
				},
				Value: 7,
				Right: &Tree{
					Left: &Tree{
						Left: nil,
						Value: 2,
						Right: nil,
					},
					Value: 12,
					Right: nil,
				},
			},
		},
		Value: 8,
		Right: &Tree{
			Left: nil,
			Value: 4,
			Right: &Tree{
				Left: &Tree{
					Left: nil,
					Value: 3,
					Right: nil,
				},
				Value: 11,
				Right: nil,
			},
		},
	}
}

func compareOrder(t *testing.T, expectedOrder []int, actualOrder []int) {
	if !reflect.DeepEqual(expectedOrder, actualOrder) {
		t.Errorf("Incorrect order:\nExpected %v\nGot %v", expectedOrder, actualOrder)
	}
}

func TestTraversePreOrder(t *testing.T) {
	tree := exampleTree()
	actualOrder := TraversePreOrder(&tree)
	expectedOrder := []int{8, 5, 9, 7, 1, 12, 2, 4, 11, 3}
	compareOrder(t, expectedOrder, actualOrder)
}

func TestTraversePostOrder(t *testing.T) {
	tree := exampleTree()
	actualOrder := TraversePostOrder(&tree)
	expectedOrder := []int{9, 1, 2, 12, 7, 5, 3, 11, 4, 8}
	compareOrder(t, expectedOrder, actualOrder)
}

func TestTraverseInOrder(t *testing.T) {
	tree := exampleTree()
	actualOrder := TraverseInOrder(&tree)
	expectedOrder := []int{9, 5, 1, 7, 2, 12, 8, 4, 3, 11}
	compareOrder(t, expectedOrder, actualOrder)
}

func TestTraverseLevelOrder(t *testing.T) {
	tree := exampleTree()
	actualOrder := TraverseLevelOrder(&tree)
	expectedOrder := []int{8, 5, 4, 9, 7, 11, 1, 12, 3, 2}
	compareOrder(t, expectedOrder, actualOrder)
}
