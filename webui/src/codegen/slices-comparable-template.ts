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
    const union = symbol('Union', true, config);
    const intersection = symbol('Intersection', true, config);
    const difference = symbol('Difference', true, config);

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
        sourceSlices := make([][]${sliceType}, 0, len(sorted))
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

    // ${difference} creates difference group of sorted slices and returns.
    func ${difference}(sorted1, sorted2 []${sliceType}) []${sliceType} {
        var result []${sliceType}
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

    // ${intersection} creates intersection group of sorted slices and returns.
    func ${intersection}(sorted ...[]${sliceType}) []${sliceType} {
        sort.Slice(sorted, func(i, j int) bool {
            return len(sorted[i]) < len(sorted[j])
        })
        var result []${sliceType}
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

    // ${union} unions sorted slices and returns new slices.
    func ${union}(sorted ...[]${sliceType}) []${sliceType} {
        length := 0
        sourceSlices := make([][]${sliceType}, 0, len(sorted))
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
        result := make([]${sliceType}, length)
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
    `;
}
