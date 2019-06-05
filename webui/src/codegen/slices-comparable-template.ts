import { GeneratorConfig, symbol } from './config';

export function generateComparable(config: GeneratorConfig): string {
    const { packageName, sliceType } = config;
    const sort = symbol('Sort', true, config);
    const binarySearch = symbol('BinarySearch', true, config);
    const indexOf = symbol('IndexOf', true, config);
    const contains = symbol('Contains', true, config);
    const insert = symbol('Insert', true, config);
    const remove = symbol('Remove', true, config);
    const iterateOver = symbol('IterateOver', true, config);
    const merge = symbol('Merge', true, config);
    const lessThan = symbol('LessThan', true, config);

    return `
    package ${packageName}

    import (
        "sort"
    )

    // ${sort} sorts an array using the provided comparator
    func ${sort}(a []${sliceType}) (err error) {
        sort.Slice(a, func(i, j int) bool {
            return a[i] < a[j]
        })
        return nil
    }

    // ${binarySearch} returns first index i that satisfies slices[i] <= item.
    func ${binarySearch}(sorted []${sliceType}, item ${sliceType}) int {
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

    // ${indexOf} returns index of item. If item is not in a sorted slice, it returns -1.
    func ${indexOf}(sorted []${sliceType}, item ${sliceType}) int {
        i := ${binarySearch}(sorted, item, lt)
        if sorted[i] == item {
            return i
        }
        return -1
    }

    // ${contains} returns true if item is in a sorted slice. Otherwise false.
    func ${contains}(sorted []${sliceType}, item ${sliceType}) bool {
        i := ${binarySearch}(sorted, item, lt)
        return sorted[i] == item
    }

    // ${insert} inserts item in correct position and returns a sorted slice.
    func ${insert}(sorted []${sliceType}, item ${sliceType}) []${sliceType} {
        i := ${binarySearch}(sorted, item, lt)
        if i == len(sorted) - 1 && sorted[i] < item {
            return append(sorted, item)
        }
        return append(sorted[:i], append([]${sliceType}{item}, sorted[i:]...)...)
    }

    // ${remove} removes item in a sorted slice.
    func ${remove}(sorted []${sliceType}, item ${sliceType}) []${sliceType} {
        i := ${binarySearch}(sorted, item, lt)
        if sorted[i] == item {
            return append(sorted[:i], sorted[i+1:]...)
        }
        return sorted
    }

    // ${iterateOver} iterates over input sorted slices and calls callback with each items in ascendant order.
    func ${iterateOver}(callback func(item ${sliceType}, srcIndex int), sorted ...[]${sliceType}) {
        sourceSlices := sorted
        sourceSliceCount := len(sorted)
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

    // ${merge} merges sorted slices and returns new slices.
    func ${merge}(sorted ...[]${sliceType}) []${sliceType} {
        length := 0
        for _, src := range sorted {
            length += len(src)
        }
        result := make([]${sliceType}, length)
        sourceSlices := sorted
        sourceSliceCount := len(sorted)
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
    `;
}
