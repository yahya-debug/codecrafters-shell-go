package main

import (
	"cmp"
)

type Pair[T, U any] struct {
	f T
	s U
}

// MergeSort is a recursive function that sorts a slice of integers using the merge sort algorithm.
func MergeSort[T cmp.Ordered](items []T) []T {
	// Base case: a slice with fewer than 2 elements is already sorted.
	if len(items) < 2 {
		return items
	}
	// Split the slice into two halves.
	mid := len(items) / 2
	left := items[:mid]
	right := items[mid:]
	// Recursively sort both halves and then merge them.
	return merge(MergeSort(left), MergeSort(right))
}

// merge combines two sorted slices into a single sorted slice.
func merge[T cmp.Ordered](left, right []T) []T {
	result := make([]T, 0, len(left)+len(right))
	i, j := 0, 0
	// Compare elements from both slices and append the smaller one to the result.
	for i < len(left) && j < len(right) {
		if left[i] < right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}
	// Append remaining elements of the left slice (if any).
	for i < len(left) {
		result = append(result, left[i])
		i++
	}
	// Append remaining elements of the right slice (if any).
	for j < len(right) {
		result = append(result, right[j])
		j++
	}
	return result
}

// Binary search
func BS[T any](arr []T, target T, l int, r int, less func(a, b T) bool) (int, bool) {
	if r < l {
		return -1, false
	}
	mid := (l + r) / 2
	if less(arr[mid], target) {
		return BS(arr, target, mid+1, r, less)
	} else if less(target, arr[mid]) {
		return BS(arr, target, l, mid-1, less)
	}
	return mid, true
}
func BSs[T any](arr []T, target T, l int, r int, less func(a, b T) bool) int {
	// fmt.Println(target)
	if r-l <= 1 {
		return l
	}
	mid := (l + r) / 2
	if less(target, arr[mid]) {
		return BSs(arr, target, l, mid, less)
	} else {
		return BSs(arr, target, mid, r, less)
	}
}

// Min & max
func minx(a *int, b int) {
	if b < *a {
		*a = b
	}
}
func maxx(a *int, b int) {
	if b > *a {
		*a = b
	}
}
