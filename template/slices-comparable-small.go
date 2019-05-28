package slices

import (
	"github.com/cheekybits/genny/generic"
	"sort"
)

type ValueType generic.Type

// ValueTypeSort sorts an array using the provided comparator
func ValueTypeSort(a []ValueType) (err error) {
	sort.Slice(a, func(i, j int) bool {
		return a[i] < a[j]
	})
	return nil
}

// ValueTypeBinarySearch returns first index i that satisfies slices[i] <= item.
func ValueTypeBinarySearch(sorted []ValueType, item ValueType) int {
	// Define f(-1) == false and f(n) == true.
	// Invariant: f(i-1) == false, f(j) == true.
	i, j := 0, len(sorted) - 1
	for i < j {
		h := int(uint(i+j) >> 1) // avoid overflow when computing h
		// i ≤ h < j
		if sorted[h] < item {
			i = h + 1 // preserves f(i-1) == false
		} else {
			j = h // preserves f(j) == true
		}
	}
	// i == j, f(i-1) == false, and f(j) (= f(i)) == true  =>  answer is i.
	return i
}

// ValueTypeIndexOf returns index of item. If item is not in a sorted slice, it returns -1.
func ValueTypeIndexOf(sorted []ValueType, item ValueType) int {
	i := ValueTypeBinarySearch(sorted, item)
	if sorted[i] == item {
		return i
	}
	return -1
}

// ValueTypeContains returns true if item is in a sorted slice. Otherwise false.
func ValueTypeContains(sorted []ValueType, item ValueType) bool {
	i := ValueTypeBinarySearch(sorted, item)
	return sorted[i] == item
}

// ValueTypeInsert inserts item in correct position and returns a sorted slice.
func ValueTypeInsert(sorted []ValueType, item ValueType) []ValueType {
	i := ValueTypeBinarySearch(sorted, item)
	if i == len(sorted) - 1 && sorted[i] < item {
		return append(sorted, item)
	}
	return append(sorted[:i], append([]ValueType{item}, sorted[i:]...)...)
}

// ValueTypeRemove removes item in a sorted slice.
func ValueTypeRemove(sorted []ValueType, item ValueType) []ValueType {
	i := ValueTypeBinarySearch(sorted, item)
	if sorted[i] == item {
		return append(sorted[:i], sorted[i+1:]...)
	}
	return sorted
}
