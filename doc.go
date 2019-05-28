// slices package provides template for genny.
//
// It generates generic algorithms for slices:
//
//   genny -in=$GOPATH/src/github.com/shibukawa/slices/template/slices.go -out=mystructslices.go gen "ValueType=MyStruct"
//
// This commands generates the following functions:
//
//    MyStructSort(slices []MyStruct, lt MyStructLessThan) error
//    MyStructBinarySearch(sorted []MyStruct, item MyStruct, lt MyStructLessThan) int
//    MyStructIndexOf(sorted []MyStruct, item MyStruct, lt MyStructLessThan) int
//    MyStructContains(sorted []MyStruct, item MyStruct, lt MyStructLessThan) bool
//    MyStructInsert(sorted []MyStruct, item MyStruct, lt MyStructLessThan) []MyStruct
//    MyStructRemove(sorted []MyStruct, item MyStruct, lt MyStructLessThan) []MyStruct
//
// To use these functions, you should define function that has signature like this:
//
//    func (a, b MyStruct) bool
//
// The function excepts MyStructSort assumes sorted slice as a first argument.
//
// ValueTypeSort
//
// This function provides timsort algorithm that is fast, stable sort algorithm. based on github.com/psilva261/timsort.
//
// ValueTypeBinarySearch
//
// This function returns first index i that satisfies slices[i] <= item.
//
// ValueTypeIndexOf
//
// This function returns index of item. If item is not in a sorted slice, it returns -1.
//
// ValueTypeContains
//
// This function returns true if item is in a sorted slice. Otherwise false.
//
// ValueTypeInsert
//
// This function inserts item in correct position and returns a sorted slice.
//
// ValueTypeRemove
//
// This function removes item in a sorted slice.
//
package slices