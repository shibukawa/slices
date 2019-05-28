# slices

genny's template to handle sorted slice.

It generates generic algorithms for slices:


```sh
$ genny -in=$GOPATH/src/github.com/shibukawa/slices/template/slices.go -out=mystructslices.go gen "ValueType=MyStruct"
```

This commands generates the following functions:

* MyStructSort(slices []MyStruct, lt MyStructLessThan) error
* MyStructBinarySearch(sorted []MyStruct, item MyStruct, lt MyStructLessThan) int
* MyStructIndexOf(sorted []MyStruct, item MyStruct, lt MyStructLessThan) int
* MyStructContains(sorted []MyStruct, item MyStruct, lt MyStructLessThan) bool
* MyStructInsert(sorted []MyStruct, item MyStruct, lt MyStructLessThan) []MyStruct
* MyStructRemove(sorted []MyStruct, item MyStruct, lt MyStructLessThan) []MyStruct

To use these functions, you should define function that compares values in slice and it has signature like this:

```go
func (a, b MyStruct) bool
```

If you want to compare field "name", the implementation should be the following code:

```go
func (a, b MyStruct) bool {
	return a.name < b.name
}
```

The function except MyStructSort assumes sorted slice as a first argument.

## Template Types

There for template files

* slices.go: Standard template. Use TimSort. Accept "LessThan" function as a comparator.
* slices-comaprable.go: Template for built-in types. Use TimSort. Use ``<`` operator as a comparator.
* slices-small.go: Standard template. Use "sort.Slices". Accept "LessThan" function as a comparator.
* slices-comaprable-small.go: Template for built-in types. Use "sort.Slices". Use ``<`` operator as a comparator.

## Generated Function Reference

### [ValueType]Sort()

This function provides timsort algorithm that is fast, stable sort algorithm. based on github.com/psilva261/timsort.
If you use template/slices-small.go or template/slices-comparable-small.go, it uses "sort.Slice" function.

### [ValueType]BinarySearch()

This function returns first index i that satisfies slices[i] <= item.

### [ValueType]IndexOf()

This function returns index of item. If item is not in a sorted slice, it returns -1.

### [ValueType]Contains()

This function returns true if item is in a sorted slice. Otherwise false.

### [ValueType]Insert()

This function insert item in correct position and returns a sorted slice.

### [ValueType]Remove()

This function remove item in a sorted slice.

## Credits/Thanks

This repository is a template for genny:

* [genny - Generics for Go](https://github.com/cheekybits/genny)

This template is highly depending on high performance sort algorithm:

* https://github.com/psilva261/timsort

This code is tested with Property based testing:

* [gopter: GOlang Property TestER](https://github.com/leanovate/gopter)

## License

Apache 2
