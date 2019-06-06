package comparablesmall

import (
	"sort"
	"testing"
	"reflect"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func TestSortInt(t *testing.T) {
	numberGenerator := gen.Int()
	numSliceGenerator := gen.SliceOf(numberGenerator)

	properties := gopter.NewProperties(nil)

	properties.Property("sort returns stable", prop.ForAll(func(input []int) bool {
		timSort := make([]int, len(input))
		defaultSort := make([]int, len(input))
		copy(timSort, input)
		copy(defaultSort, input)

		IntSort(timSort)
		sort.Ints(defaultSort)
		return reflect.DeepEqual(timSort, defaultSort)
	}, numSliceGenerator))

	properties.TestingRun(t)
}

func TestBinarySearch(t *testing.T) {
	numberGenerator := gen.Int()
	numSliceGenerator := gen.SliceOfN(20, numberGenerator)

	properties := gopter.NewProperties(nil)

	properties.Property("binary search found items", prop.ForAll(func(input []int) bool {
		value := input[0]
		orig := make([]int, len(input))
		copy(orig, input)
		IntSort(input)
		i := IntBinarySearch(input, value)
		if input[i] != value {
			t.Log(i, value, input[i])
			t.Log(orig)
			t.Log(input)
		}
		return input[i] == value
	}, numSliceGenerator))

	properties.TestingRun(t)
}

func TestIndexOf(t *testing.T) {
	numberGenerator := gen.Int()
	numSliceGenerator := gen.SliceOfN(20, numberGenerator)

	properties := gopter.NewProperties(nil)

	properties.Property("indexOf found items", prop.ForAll(func(input []int) bool {
		value := input[0]
		IntSort(input)
		i := IntIndexOf(input, value)
		return i != -1 && input[i] == value
	}, numSliceGenerator))

	properties.Property("indexOf returns -1 if not found", prop.ForAll(func(input []int) bool {
		value := input[0]
		array := input[1:]
		IntSort(array)
		i := IntIndexOf(array, value)
		return i == -1
	}, numSliceGenerator))

	properties.TestingRun(t)
}

func TestContains(t *testing.T) {
	numberGenerator := gen.Int()
	numSliceGenerator := gen.SliceOfN(20, numberGenerator)

	properties := gopter.NewProperties(nil)

	properties.Property("contains returns true if item found", prop.ForAll(func(input []int) bool {
		value := input[0]
		IntSort(input)
		return IntContains(input, value)
	}, numSliceGenerator))

	properties.Property("indexOf returns false if not found", prop.ForAll(func(input []int) bool {
		value := input[0]
		array := input[1:]
		IntSort(array)
		return !IntContains(array, value)
	}, numSliceGenerator))

	properties.TestingRun(t)
}

func TestInsert(t *testing.T) {
	numberGenerator := gen.Int()
	numSliceGenerator := gen.SliceOfN(20, numberGenerator)

	properties := gopter.NewProperties(nil)

	properties.Property("insert returns new sorted slices", prop.ForAll(func(input []int) bool {
		expected := make([]int, len(input))
		copy(expected, input)
		IntSort(expected)

		value := input[0]
		array := input[1:]
		IntSort(array)

		inserted := IntInsert(array, value)

		return reflect.DeepEqual(expected, inserted)
	}, numSliceGenerator))

	properties.TestingRun(t)
}

func TestRemove(t *testing.T) {
	numberGenerator := gen.Int()
	numSliceGenerator := gen.SliceOfN(20, numberGenerator)

	properties := gopter.NewProperties(nil)

	properties.Property("remove removes item of array", prop.ForAll(func(input []int) bool {
		value := input[0]
		IntSort(input)

		removedArray := IntRemove(input, value)

		return len(removedArray) == len(input) -1 && !IntContains(removedArray, value)
	}, numSliceGenerator))

	properties.TestingRun(t)
}

func TestMerge(t *testing.T) {
	numberGenerator := gen.Int()
	numSliceGenerator := gen.SliceOf(numberGenerator)

	properties := gopter.NewProperties(nil)

	properties.Property("marge item of slices", prop.ForAll(func(input1, input2, input3 []int) bool {
		IntSort(input1)
		IntSort(input2)
		IntSort(input3)

		marged := IntMerge(input1, input2, input3)

		if len(marged) == 0 {
			return true
		}

		expected := make([]int, len(marged))
		copy(expected, marged)
		IntSort(expected)

		return reflect.DeepEqual(expected, marged)
	}, numSliceGenerator, numSliceGenerator, numSliceGenerator))

	properties.TestingRun(t)
}

func TestIterateOver(t *testing.T) {
	numberGenerator := gen.Int()
	numSliceGenerator := gen.SliceOf(numberGenerator)

	properties := gopter.NewProperties(nil)

	properties.Property("iterate item of slices", prop.ForAll(func(input1, input2, input3 []int) bool {
		IntSort(input1)
		IntSort(input2)
		IntSort(input3)

		var result []int
		IntIterateOver(func(item, srcIndex int) {
			result = append(result, item)
		}, input1, input2, input3)

		if len(result) == 0 {
			return true
		}

		expected := make([]int, len(result))
		copy(expected, result)
		IntSort(expected)

		return reflect.DeepEqual(expected, result)
	}, numSliceGenerator, numSliceGenerator, numSliceGenerator))

	properties.TestingRun(t)
}