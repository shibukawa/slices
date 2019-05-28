// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package small

import "sort"

// IntLessThan is Delegate type that sorting uses as a comparator
type IntLessThan func(a, b int) bool

// IntSort sorts an array using the provided comparator
func IntSort(a []int, lt IntLessThan) (err error) {
	sort.Slice(a, func(i, j int) bool {
		return lt(a[i], a[j])
	})
	return nil
}

// IntBinarySearch returns first index i that satisfies slices[i] <= item.
func IntBinarySearch(sorted []int, item int, lt IntLessThan) int {
	// Define f(-1) == false and f(n) == true.
	// Invariant: f(i-1) == false, f(j) == true.
	i, j := 0, len(sorted)-1
	for i < j {
		h := int(uint(i+j) >> 1) // avoid overflow when computing h
		// i ≤ h < j
		if lt(sorted[h], item) {
			i = h + 1 // preserves f(i-1) == false
		} else {
			j = h // preserves f(j) == true
		}
	}
	// i == j, f(i-1) == false, and f(j) (= f(i)) == true  =>  answer is i.
	return i
}

// IntIndexOf returns index of item. If item is not in a sorted slice, it returns -1.
func IntIndexOf(sorted []int, item int, lt IntLessThan) int {
	i := IntBinarySearch(sorted, item, lt)
	if !lt(sorted[i], item) && !lt(item, sorted[i]) {
		return i
	}
	return -1
}

// IntContains returns true if item is in a sorted slice. Otherwise false.
func IntContains(sorted []int, item int, lt IntLessThan) bool {
	i := IntBinarySearch(sorted, item, lt)
	return !lt(sorted[i], item) && !lt(item, sorted[i])
}

// IntInsert inserts item in correct position and returns a sorted slice.
func IntInsert(sorted []int, item int, lt IntLessThan) []int {
	i := IntBinarySearch(sorted, item, lt)
	if i == len(sorted)-1 && lt(sorted[i], item) {
		return append(sorted, item)
	}
	return append(sorted[:i], append([]int{item}, sorted[i:]...)...)
}

// IntRemove removes item in a sorted slice.
func IntRemove(sorted []int, item int, lt IntLessThan) []int {
	i := IntBinarySearch(sorted, item, lt)
	if !lt(sorted[i], item) && !lt(item, sorted[i]) {
		return append(sorted[:i], sorted[i+1:]...)
	}
	return sorted
}
