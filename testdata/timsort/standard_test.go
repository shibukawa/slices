package standard

import (
	"sort"
	"testing"
	"reflect"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func cmp(a, b int) bool {
	return a < b
}

func TestSortInt(t *testing.T) {
	numberGenerator := gen.Int()
	numSliceGenerator := gen.SliceOf(numberGenerator)

	properties := gopter.NewProperties(nil)

	properties.Property("sort returns stable", prop.ForAll(func(input []int) bool {
		timSort := make([]int, len(input))
		defaultSort := make([]int, len(input))
		copy(timSort, input)
		copy(defaultSort, input)

		IntSort(timSort, cmp)
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
		IntSort(input, cmp)
		i := IntBinarySearch(input, value, cmp)
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
		IntSort(input, cmp)
		i := IntIndexOf(input, value, cmp)
		return i != -1 && input[i] == value
	}, numSliceGenerator))

	properties.Property("indexOf returns -1 if not found", prop.ForAll(func(input []int) bool {
		value := input[0]
		array := input[1:]
		IntSort(array, cmp)
		i := IntIndexOf(array, value, cmp)
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
		IntSort(input, cmp)
		return IntContains(input, value, cmp)
	}, numSliceGenerator))

	properties.Property("indexOf returns false if not found", prop.ForAll(func(input []int) bool {
		value := input[0]
		array := input[1:]
		IntSort(array, cmp)
		return !IntContains(array, value, cmp)
	}, numSliceGenerator))

	properties.TestingRun(t)
}

func TestInsert(t *testing.T) {
	numberGenerator := gen.Int()
	numSliceGenerator := gen.SliceOfN(20, numberGenerator)

	properties := gopter.NewProperties(nil)

	properties.Property("insert returns new sorted slice", prop.ForAll(func(input []int) bool {
		expected := make([]int, len(input))
		copy(expected, input)
		IntSort(expected, cmp)

		value := input[0]
		array := input[1:]
		IntSort(array, cmp)

		inserted := IntInsert(array, value, cmp)

		return reflect.DeepEqual(expected, inserted)
	}, numSliceGenerator))

	properties.TestingRun(t)
}

func TestRemove(t *testing.T) {
	numberGenerator := gen.Int()
	numSliceGenerator := gen.SliceOfN(20, numberGenerator)

	properties := gopter.NewProperties(nil)

	properties.Property("remove removes item of slice", prop.ForAll(func(input []int) bool {
		value := input[0]
		IntSort(input, cmp)

		removedArray := IntRemove(input, value, cmp)

		return len(removedArray) == len(input) -1 && !IntContains(removedArray, value, cmp)
	}, numSliceGenerator))

	properties.TestingRun(t)
}

func TestUnion(t *testing.T) {
	numberGenerator := gen.Int()
	numSliceGenerator := gen.SliceOf(numberGenerator)

	properties := gopter.NewProperties(nil)

	properties.Property("union item of slices", prop.ForAll(func(input1, input2, input3 []int) bool {
		IntSort(input1, cmp)
		IntSort(input2, cmp)
		IntSort(input3, cmp)

		union := IntUnion(cmp, input1, input2, input3)

		if len(union) == 0 {
			return true
		}

		expected := make([]int, len(union))
		copy(expected, union)
		IntSort(expected, cmp)

		return reflect.DeepEqual(expected, union)
	}, numSliceGenerator, numSliceGenerator, numSliceGenerator))

	properties.TestingRun(t)
}

func TestIterateOver(t *testing.T) {
	numberGenerator := gen.Int()
	numSliceGenerator := gen.SliceOf(numberGenerator)

	properties := gopter.NewProperties(nil)

	properties.Property("iterate item of slices", prop.ForAll(func(input1, input2, input3 []int) bool {
		IntSort(input1, cmp)
		IntSort(input2, cmp)
		IntSort(input3, cmp)

		var result []int
		IntIterateOver(cmp, func(item, srcIndex int) {
			result = append(result, item)
		}, input1, input2, input3)

		if len(result) == 0 {
			return true
		}

		expected := make([]int, len(result))
		copy(expected, result)
		IntSort(expected, cmp)
		return reflect.DeepEqual(expected, result)
	}, numSliceGenerator, numSliceGenerator, numSliceGenerator))

	properties.TestingRun(t)
}