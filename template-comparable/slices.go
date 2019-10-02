package template_comparable

import (
	"github.com/cheekybits/genny/generic"
	"sort"
)

type ValueType generic.Number

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
	i, j := 0, len(sorted)-1
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
	if i == len(sorted)-1 && sorted[i] < item {
		return append(sorted, item)
	}
	return append(sorted[:i], append([]ValueType{item}, sorted[i:]...)...)
}

// ValueTypeRemove removes item in a sorted slice.
func ValueTypeRemove(sorted []ValueType, item ValueType) []ValueType {
	i := ValueTypeBinarySearch(sorted, item)
	if sorted[i] == item {
		return ValueTypeRemoveAt(sorted, i)
	}
	return sorted
}

// ValueTypeRemoveAt removes item in a slice.
func ValueTypeRemoveAt(sorted []ValueType, i int) []ValueType {
	return append(sorted[:i], sorted[i+1:]...)
}

// ValueTypeIterateOver iterates over input sorted slices and calls callback with each items in ascendant order.
func ValueTypeIterateOver(callback func(item ValueType, srcIndex int), sorted ...[]ValueType) {
	sourceSlices := make([][]ValueType, 0, len(sorted))
	for _, src := range sorted {
		if len(src) > 0 {
			sourceSlices = append(sourceSlices, src)
		}
	}
	sourceSliceCount := len(sourceSlices)
	if sourceSliceCount == 0 {
		return
	} else if sourceSliceCount == 1 {
		for i, value := range sourceSlices[0] {
			callback(value, i)
		}
		return
	}
	indexes := make([]int, sourceSliceCount)
	sliceIndex := make([]int, sourceSliceCount)
	for i := range sourceSlices {
		sliceIndex[i] = i
	}
	index := 0
	for {
		minSlice := 0
		minItem := sourceSlices[0][indexes[0]]
		for i := 1; i < sourceSliceCount; i++ {
			if sourceSlices[i][indexes[i]] < minItem {
				minSlice = i
				minItem = sourceSlices[i][indexes[i]]
			}
		}
		callback(minItem, sliceIndex[minSlice])
		index++
		indexes[minSlice]++
		if indexes[minSlice] == len(sourceSlices[minSlice]) {
			sourceSlices = append(sourceSlices[:minSlice], sourceSlices[minSlice+1:]...)
			indexes = append(indexes[:minSlice], indexes[minSlice+1:]...)
			sliceIndex = append(sliceIndex[:minSlice], sliceIndex[minSlice+1:]...)
			sourceSliceCount--
			if len(sourceSlices) == 1 {
				slice := sourceSlices[0]
				for i := indexes[0]; i < len(slice); i++ {
					callback(slice[i], sliceIndex[0])
				}
				return
			}
		}
	}
}

// ValueTypeUnion unions sorted slices and returns new slices.
func ValueTypeUnion(sorted ...[]ValueType) []ValueType {
	length := 0
	sourceSlices := make([][]ValueType, 0, len(sorted))
	for _, src := range sorted {
		if len(src) > 0 {
			length += len(src)
			sourceSlices = append(sourceSlices, src)
		}
	}
	if length == 0 {
		return nil
	} else if len(sourceSlices) == 1 {
		return sourceSlices[0]
	}
	result := make([]ValueType, length)
	sourceSliceCount := len(sourceSlices)
	indexes := make([]int, sourceSliceCount)
	index := 0
	for {
		minSlice := 0
		minItem := sourceSlices[0][indexes[0]]
		for i := 1; i < sourceSliceCount; i++ {
			if sourceSlices[i][indexes[i]] < minItem {
				minSlice = i
				minItem = sourceSlices[i][indexes[i]]
			}
		}
		result[index] = minItem
		index++
		indexes[minSlice]++
		if indexes[minSlice] == len(sourceSlices[minSlice]) {
			sourceSlices = append(sourceSlices[:minSlice], sourceSlices[minSlice+1:]...)
			indexes = append(indexes[:minSlice], indexes[minSlice+1:]...)
			sourceSliceCount--
			if len(sourceSlices) == 1 {
				copy(result[index:], sourceSlices[0][indexes[0]:])
				return result
			}
		}
	}
}

func ValueTypeDifference(sorted1, sorted2 []ValueType) []ValueType {
	var result []ValueType
	var i, j int
	for i < len(sorted1) && j < len(sorted2) {
		if sorted1[i] < sorted2[j] {
			result = append(result, sorted1[i])
			i++
		} else if sorted2[j] < sorted1[i] {
			j++
		} else {
			i++
			j++
		}
	}
	result = append(result, sorted1[i:]...)
	return result
}

func ValueTypeIntersection(sorted ...[]ValueType) []ValueType {
	sort.Slice(sorted, func(i, j int) bool {
		return len(sorted[i]) < len(sorted[j])
	})
	var result []ValueType
	if len(sorted[0]) == 0 {
		return result
	}
	cursors := make([]int, len(sorted))
	terminate := false
	for _, value := range sorted[0] {
		needIncrement := false
		for i := 1; i < len(sorted); i++ {
			found := false
			for j := cursors[i]; j < len(sorted[i]); j++ {
				valueOfOtherSlice := sorted[i][cursors[i]]
				if valueOfOtherSlice < value {
					cursors[i] = j + 1
				} else if value < valueOfOtherSlice {
					needIncrement = true
					break
				} else {
					found = true
					break
				}
			}
			if needIncrement {
				break
			}
			if !found {
				terminate = true
				break
			}
		}
		if terminate {
			break
		}
		if !needIncrement {
			result = append(result, value)
		}
	}
	return result
}
