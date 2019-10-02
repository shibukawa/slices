import { GeneratorConfig, symbol } from './config';

export function generateComparableTimsort(config: GeneratorConfig): string {
    const { packageName, sliceType } = config;
    const sort = symbol('Sort', true, config);
    const binarySearch = symbol('BinarySearch', true, config);
    const indexOf = symbol('IndexOf', true, config);
    const contains = symbol('Contains', true, config);
    const insert = symbol('Insert', true, config);
    const remove = symbol('Remove', true, config);
    const iterateOver = symbol('IterateOver', true, config);
    const union = symbol('Union', true, config);

    const countRunAndMakeAscending = symbol('countRunAndMakeAscending', false, config);
    const reverseRange = symbol('reverseRange', false, config);
    const minRunLength = symbol('minRunLength', false, config);
    const iBinarySort = symbol('binarySort', false, config);
    const timSortHandler = symbol('timSortHandler', false, config);
    const gallopLeft = symbol('gallopLeft', false, config);
    const gallopRight = symbol('gallopRight', false, config);

    return `
    package ${packageName}

    import (
        "errors"
        "fmt"
    )

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
        i := ${binarySearch}(sorted, item)
        if sorted[i] == item {
            return i
        }
        return -1
    }

    // ${contains} returns true if item is in a sorted slice. Otherwise false.
    func ${contains}(sorted []${sliceType}, item ${sliceType}) bool {
        i := ${binarySearch}(sorted, item)
        return sorted[i] == item
    }

    // ${insert} inserts item in correct position and returns a sorted slice.
    func ${insert}(sorted []${sliceType}, item ${sliceType}) []${sliceType} {
        i := ${binarySearch}(sorted, item)
        if i == len(sorted) - 1 && sorted[i] < item {
            return append(sorted, item)
        }
        return append(sorted[:i], append([]${sliceType}{item}, sorted[i:]...)...)
    }

    // ${remove} removes item in a sorted slice.
    func ${remove}(sorted []${sliceType}, item ${sliceType}) []${sliceType} {
        i := ${binarySearch}(sorted, item)
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

    // timsort code is based on https://github.com/psilva261/timsort
    //
    // Package timsort provides fast stable sort, uses external comparator.
    //
    // A stable, adaptive, iterative mergesort that requires far fewer than
    // n lg(n) comparisons when running on partially sorted arrays, while
    // offering performance comparable to a traditional mergesort when run
    // on random arrays.  Like all proper mergesorts, this sort is stable and
    // runs O(n log n) time (worst case).  In the worst case, this sort requires
    // temporary storage space for n/2 object references; in the best case,
    // it requires only a small constant amount of space.
    //
    // This implementation was derived from Java's TimSort object by Josh Bloch,
    // which, in turn, was based on the original code by Tim Peters:
    //
    // http://svn.python.org/projects/python/trunk/Objects/listsort.txt
    //
    // Mike K.

    // ${sort} sorts an array using the provided comparator
    func ${sort}(a []${sliceType}) (err error) {
        minMerge := 32
        lo := 0
        hi := len(a)
        nRemaining := hi

        if nRemaining < 2 {
            return // Arrays of size 0 and 1 are always sorted
        }

        if nRemaining < minMerge {
            initRunLen, err := ${countRunAndMakeAscending}(a, lo, hi)
            if err != nil {
                return err
            }
            return ${iBinarySort}(a, lo, hi, lo+initRunLen)
        }
        ts := ${symbol('newTimSort', false, config)}(a)
        minRun, err := ${minRunLength}(nRemaining)
        if err != nil {
            return
        }
        for {
            // Identify next run
            runLen, err := ${countRunAndMakeAscending}(a, lo, hi)
            if err != nil {
                return err
            }
            if runLen < minRun {
                force := minRun
                if nRemaining <= minRun {
                    force = nRemaining
                }
                if err = ${iBinarySort}(a, lo, lo+force, lo+runLen); err != nil {
                    return err
                }
                runLen = force
            }
            ts.pushRun(lo, runLen)
            if err = ts.mergeCollapse(); err != nil {
                return err
            }
            lo += runLen
            nRemaining -= runLen
            if nRemaining == 0 {
                break
            }
        }
        if lo != hi {
            return errors.New("lo must equal hi")
        }
        if err = ts.mergeForceCollapse(); err != nil {
            return
        }
        if ts.stackSize != 1 {
            return errors.New("ts.stackSize != 1")
        }
        return
    }

    type ${timSortHandler} struct {
        a []${sliceType}
        minGallop int
        tmp []${sliceType}
        stackSize int
        runBase   []int
        runLen    []int
    }

    func ${symbol('newTimSort', false, config)}(a []${sliceType}) (h *${timSortHandler}) {
        minGallop := 7
        initialTmpStorageLength := 256
        h = new(${timSortHandler})

        h.a = a
        h.minGallop = minGallop
        h.stackSize = 0
        len := len(a)
        tmpSize := initialTmpStorageLength
        if len < 2*tmpSize {
            tmpSize = len / 2
        }
        h.tmp = make([]${sliceType}, tmpSize)
        stackLen := 40
        if len < 120 {
            stackLen = 5
        } else if len < 1542 {
            stackLen = 10
        } else if len < 119151 {
            stackLen = 19
        }
        h.runBase = make([]int, stackLen)
        h.runLen = make([]int, stackLen)
        return h
    }

    func ${iBinarySort}(a []${sliceType}, lo, hi, start int) (err error) {
        if lo > start || start > hi {
            return errors.New("lo <= start && start <= hi")
        }
        if start == lo {
            start++
        }
        for ; start < hi; start++ {
            pivot := a[start]
            left := lo
            right := start
            if left > right {
                return errors.New("left <= right")
            }
            for left < right {
                mid := int(uint(left+right) >> 1)
                if pivot < a[mid] {
                    right = mid
                } else {
                    left = mid + 1
                }
            }
            if left != right {
                return errors.New("left == right")
            }
            n := start - left // The number of elements to move
            // just an optimization for copy in default case
            if n <= 2 {
                if n == 2 {
                    a[left+2] = a[left+1]
                }
                if n > 0 {
                    a[left+1] = a[left]
                }
            } else {
                copy(a[left+1:], a[left:left+n])
            }
            a[left] = pivot
        }
        return
    }

    func ${countRunAndMakeAscending}(a []${sliceType}, lo, hi int) (int, error) {
        if lo >= hi {
            return 0, errors.New("lo < hi")
        }
        runHi := lo + 1
        if runHi == hi {
            return 1, nil
        }
        if a[runHi] < a[lo] { // Descending
            runHi++
            for runHi < hi && a[runHi] < a[runHi-1] {
                runHi++
            }
            ${reverseRange}(a, lo, runHi)
        } else { // Ascending
            for runHi < hi && !(a[runHi] < a[runHi-1]) {
                runHi++
            }
        }
        return runHi - lo, nil
    }

    func ${reverseRange}(a []${sliceType}, lo, hi int) {
        hi--
        for lo < hi {
            a[lo], a[hi] = a[hi], a[lo]
            lo++
            hi--
        }
    }

    func ${minRunLength}(n int) (int, error) {
        minMerge := 32
        if n < 0 {
            return 0, errors.New("n >= 0")
        }
        r := 0 // Becomes 1 if any 1 bits are shifted off
        for n >= minMerge {
            r |= (n & 1)
            n >>= 1
        }
        return n + r, nil
    }

    func (h *${timSortHandler}) pushRun(runBase, runLen int) {
        h.runBase[h.stackSize] = runBase
        h.runLen[h.stackSize] = runLen
        h.stackSize++
    }

    func (h *${timSortHandler}) mergeCollapse() (err error) {
        for h.stackSize > 1 {
            n := h.stackSize - 2
            if (n > 0 && h.runLen[n-1] <= h.runLen[n]+h.runLen[n+1]) ||
                (n > 1 && h.runLen[n-2] <= h.runLen[n-1]+h.runLen[n]) {
                if h.runLen[n-1] < h.runLen[n+1] {
                    n--
                }
                if err = h.mergeAt(n); err != nil {
                    return
                }
            } else if h.runLen[n] <= h.runLen[n+1] {
                if err = h.mergeAt(n); err != nil {
                    return
                }
            } else {
                break // Invariant is established
            }
        }
        return
    }

    func (h *${timSortHandler}) mergeForceCollapse() (err error) {
        for h.stackSize > 1 {
            n := h.stackSize - 2
            if n > 0 && h.runLen[n-1] < h.runLen[n+1] {
                n--
            }
            if err = h.mergeAt(n); err != nil {
                return
            }
        }
        return
    }

    func (h *${timSortHandler}) mergeAt(i int) (err error) {
        if h.stackSize < 2 {
            return errors.New("stackSize >= 2")
        }

        if i < 0 {
            return errors.New(" i >= 0")
        }

        if i != h.stackSize-2 && i != h.stackSize-3 {
            return errors.New("if i == stackSize - 2 || i == stackSize - 3")
        }

        base1 := h.runBase[i]
        len1 := h.runLen[i]
        base2 := h.runBase[i+1]
        len2 := h.runLen[i+1]

        if len1 <= 0 || len2 <= 0 {
            return errors.New("len1 > 0 && len2 > 0")
        }

        if base1+len1 != base2 {
            return errors.New("base1 + len1 == base2")
        }

        h.runLen[i] = len1 + len2
        if i == h.stackSize-3 {
            h.runBase[i+1] = h.runBase[i+2]
            h.runLen[i+1] = h.runLen[i+2]
        }
        h.stackSize--

        k, err := ${gallopRight}(h.a[base2], h.a, base1, len1, 0)
        if err != nil {
            return err
        }
        if k < 0 {
            return errors.New(" k >= 0;")
        }
        base1 += k
        len1 -= k
        if len1 == 0 {
            return
        }

        len2, err = ${gallopLeft}(h.a[base1+len1-1], h.a, base2, len2, len2-1)
        if err != nil {
            return
        }
        if len2 < 0 {
            return errors.New(" len2 >= 0;")
        }
        if len2 == 0 {
            return
        }

        if len1 <= len2 {
            err = h.mergeLo(base1, len1, base2, len2)
            if err != nil {
                return fmt.Errorf("mergeLo: %v", err)
            }
        } else {
            err = h.mergeHi(base1, len1, base2, len2)
            if err != nil {
                return fmt.Errorf("mergeHi: %v", err)
            }
        }
        return
    }

    func ${gallopLeft}(key ${sliceType}, a []${sliceType}, base, len, hint int) (int, error) {
        if len <= 0 || hint < 0 || hint >= len {
            return 0, errors.New(" len > 0 && hint >= 0 && hint < len;")
        }
        lastOfs := 0
        ofs := 1

        if a[base+hint] < key {
            // Gallop right until a[base+hint+lastOfs] < key <= a[base+hint+ofs]
            maxOfs := len - hint
            for ofs < maxOfs && a[base+hint+ofs] < key {
                lastOfs = ofs
                ofs = (ofs << 1) + 1
                if ofs <= 0 { // int overflow
                    ofs = maxOfs
                }
            }
            if ofs > maxOfs {
                ofs = maxOfs
            }
            lastOfs += hint
            ofs += hint
        } else { // key <= a[base + hint]
            // Gallop left until a[base+hint-ofs] < key <= a[base+hint-lastOfs]
            maxOfs := hint + 1
            for ofs < maxOfs && !(a[base+hint-ofs] < key) {
                lastOfs = ofs
                ofs = (ofs << 1) + 1
                if ofs <= 0 { // int overflow
                    ofs = maxOfs
                }
            }
            if ofs > maxOfs {
                ofs = maxOfs
            }
            tmp := lastOfs
            lastOfs = hint - ofs
            ofs = hint - tmp
        }

        if -1 > lastOfs || lastOfs >= ofs || ofs > len {
            return 0, errors.New(" -1 <= lastOfs && lastOfs < ofs && ofs <= len;")
        }

        lastOfs++
        for lastOfs < ofs {
            m := lastOfs + (ofs-lastOfs)/2
            if a[base+m] < key {
                lastOfs = m + 1 // a[base + m] < key
            } else {
                ofs = m // key <= a[base + m]
            }
        }

        if lastOfs != ofs {
            return 0, errors.New(" lastOfs == ofs") // so a[base + ofs - 1] < key <= a[base + ofs]
        }
        return ofs, nil
    }

    func ${gallopRight}(key ${sliceType}, a []${sliceType}, base, len, hint int) (int, error) {
        if len <= 0 || hint < 0 || hint >= len {
            return 0, errors.New(" len > 0 && hint >= 0 && hint < len;")
        }
        ofs := 1
        lastOfs := 0
        if key < a[base+hint] {
            // Gallop left until a[b+hint - ofs] <= key < a[b+hint - lastOfs]
            maxOfs := hint + 1
            for ofs < maxOfs && key < a[base+hint-ofs] {
                lastOfs = ofs
                ofs = (ofs << 1) + 1
                if ofs <= 0 { // int overflow
                    ofs = maxOfs
                }
            }
            if ofs > maxOfs {
                ofs = maxOfs
            }

            tmp := lastOfs
            lastOfs = hint - ofs
            ofs = hint - tmp
        } else { // a[b + hint] <= key
            // Gallop right until a[b+hint + lastOfs] <= key < a[b+hint + ofs]
            maxOfs := len - hint
            for ofs < maxOfs && !(key < a[base+hint+ofs]) {
                lastOfs = ofs
                ofs = (ofs << 1) + 1
                if ofs <= 0 { // int overflow
                    ofs = maxOfs
                }
            }
            if ofs > maxOfs {
                ofs = maxOfs
            }

            lastOfs += hint
            ofs += hint
        }
        if -1 > lastOfs || lastOfs >= ofs || ofs > len {
            return 0, errors.New("-1 <= lastOfs && lastOfs < ofs && ofs <= len")
        }

        lastOfs++
        for lastOfs < ofs {
            m := lastOfs + (ofs-lastOfs)/2

            if key < a[base+m] {
                ofs = m // key < a[b + m]
            } else {
                lastOfs = m + 1 // a[b + m] <= key
            }
        }
        if lastOfs != ofs {
            return 0, errors.New(" lastOfs == ofs") // so a[b + ofs - 1] <= key < a[b + ofs]
        }
        return ofs, nil
    }

    func (h *${timSortHandler}) mergeLo(base1, len1, base2, len2 int) (err error) {
        if len1 <= 0 || len2 <= 0 || base1+len1 != base2 {
            return errors.New(" len1 > 0 && len2 > 0 && base1 + len1 == base2")
        }

        a := h.a // For performance
        tmp := h.ensureCapacity(len1)

        copy(tmp, a[base1:base1+len1])

        cursor1 := 0     // Indexes into tmp array
        cursor2 := base2 // Indexes int a
        dest := base1    // Indexes int a

        a[dest] = a[cursor2]
        dest++
        cursor2++
        len2--
        if len2 == 0 {
            copy(a[dest:dest+len1], tmp)
            return
        }
        if len1 == 1 {
            copy(a[dest:dest+len2], a[cursor2:cursor2+len2])
            a[dest+len2] = tmp[cursor1] // Last elt of run 1 to end of merge
            return
        }

        minGallop := h.minGallop

    outer:
        for {
            count1 := 0 // Number of times in a row that first run won
            count2 := 0 // Number of times in a row that second run won

            for {
                if len1 <= 1 || len2 <= 0 {
                    return errors.New(" len1 > 1 && len2 > 0")
                }

                if a[cursor2] < tmp[cursor1] {
                    a[dest] = a[cursor2]
                    dest++
                    cursor2++
                    count2++
                    count1 = 0
                    len2--
                    if len2 == 0 {
                        break outer
                    }
                } else {
                    a[dest] = tmp[cursor1]
                    dest++
                    cursor1++
                    count1++
                    count2 = 0
                    len1--
                    if len1 == 1 {
                        break outer
                    }
                }
                if (count1 | count2) >= minGallop {
                    break
                }
            }

            for {
                if len1 <= 1 || len2 <= 0 {
                    return errors.New("len1 > 1 && len2 > 0")
                }
                count1, err = ${gallopRight}(a[cursor2], tmp, cursor1, len1, 0)
                if err != nil {
                    return
                }
                if count1 != 0 {
                    copy(a[dest:dest+count1], tmp[cursor1:cursor1+count1])
                    dest += count1
                    cursor1 += count1
                    len1 -= count1
                    if len1 <= 1 { // len1 == 1 || len1 == 0
                        break outer
                    }
                }
                a[dest] = a[cursor2]
                dest++
                cursor2++
                len2--
                if len2 == 0 {
                    break outer
                }

                count2, err = ${gallopLeft}(tmp[cursor1], a, cursor2, len2, 0)
                if err != nil {
                    return
                }
                if count2 != 0 {
                    copy(a[dest:dest+count2], a[cursor2:cursor2+count2])
                    dest += count2
                    cursor2 += count2
                    len2 -= count2
                    if len2 == 0 {
                        break outer
                    }
                }
                a[dest] = tmp[cursor1]
                dest++
                cursor1++
                len1--
                if len1 == 1 {
                    break outer
                }
                minGallop--
                if count1 < minGallop && count2 < minGallop {
                    break
                }
            }
            if minGallop < 0 {
                minGallop = 0
            }
            minGallop += 2 // Penalize for leaving gallop mode
        }

        if minGallop < 1 {
            minGallop = 1
        }
        h.minGallop = minGallop // Write back to field

        if len1 == 1 {
            if len2 <= 0 {
                return errors.New(" len2 > 0;")
            }
            copy(a[dest:dest+len2], a[cursor2:cursor2+len2])
            a[dest+len2] = tmp[cursor1] //  Last elt of run 1 to end of merge
        } else if len1 == 0 {
            return errors.New("comparison method violates its general contract")
        } else {
            if len2 != 0 {
                return errors.New("len2 == 0;")
            }
            if len1 <= 1 {
                return errors.New(" len1 > 1;")
            }

            copy(a[dest:dest+len1], tmp[cursor1:cursor1+len1])
        }
        return
    }

    func (h *${timSortHandler}) mergeHi(base1, len1, base2, len2 int) (err error) {
        if len1 <= 0 || len2 <= 0 || base1+len1 != base2 {
            return errors.New("len1 > 0 && len2 > 0 && base1 + len1 == base2;")
        }

        a := h.a // For performance
        tmp := h.ensureCapacity(len2)

        copy(tmp, a[base2:base2+len2])

        cursor1 := base1 + len1 - 1 // Indexes into a
        cursor2 := len2 - 1         // Indexes into tmp array
        dest := base2 + len2 - 1    // Indexes into a

        a[dest] = a[cursor1]
        dest--
        cursor1--
        len1--
        if len1 == 0 {
            dest -= len2 - 1
            copy(a[dest:dest+len2], tmp)
            return
        }
        if len2 == 1 {
            dest -= len1 - 1
            cursor1 -= len1 - 1
            copy(a[dest:dest+len1], a[cursor1:cursor1+len1])
            a[dest-1] = tmp[cursor2]
            return
        }
        minGallop := h.minGallop
    outer:
        for {
            count1 := 0 // Number of times in a row that first run won
            count2 := 0 // Number of times in a row that second run won

            for {
                if len1 <= 0 || len2 <= 1 {
                    return errors.New(" len1 > 0 && len2 > 1;")
                }
                if tmp[cursor2] < a[cursor1] {
                    a[dest] = a[cursor1]
                    dest--
                    cursor1--
                    count1++
                    count2 = 0
                    len1--
                    if len1 == 0 {
                        break outer
                    }
                } else {
                    a[dest] = tmp[cursor2]
                    dest--
                    cursor2--
                    count2++
                    count1 = 0
                    len2--
                    if len2 == 1 {
                        break outer
                    }
                }
                if (count1 | count2) >= minGallop {
                    break
                }
            }

            for {
                if len1 <= 0 || len2 <= 1 {
                    return errors.New(" len1 > 0 && len2 > 1;")
                }
                if gr, err := ${gallopRight}(tmp[cursor2], a, base1, len1, len1-1); err == nil {
                    count1 = len1 - gr
                } else {
                    return err
                }
                if count1 != 0 {
                    dest -= count1
                    cursor1 -= count1
                    len1 -= count1
                    copy(a[dest+1:dest+1+count1], a[cursor1+1:cursor1+1+count1])
                    if len1 == 0 {
                        break outer
                    }
                }
                a[dest] = tmp[cursor2]
                dest--
                cursor2--
                len2--
                if len2 == 1 {
                    break outer
                }

                if gl, err := ${gallopLeft}(a[cursor1], tmp, 0, len2, len2-1); err == nil {
                    count2 = len2 - gl
                } else {
                    return err
                }
                if count2 != 0 {
                    dest -= count2
                    cursor2 -= count2
                    len2 -= count2
                    copy(a[dest+1:dest+1+count2], tmp[cursor2+1:cursor2+1+count2])
                    if len2 <= 1 { // len2 == 1 || len2 == 0
                        break outer
                    }
                }
                a[dest] = a[cursor1]
                dest--
                cursor1--
                len1--
                if len1 == 0 {
                    break outer
                }
                minGallop--

                if count1 < minGallop && count2 < minGallop {
                    break
                }
            }
            if minGallop < 0 {
                minGallop = 0
            }
            minGallop += 2 // Penalize for leaving gallop mode
        } // End of "outer" loop

        if minGallop < 1 {
            minGallop = 1
        }

        h.minGallop = minGallop // Write back to field

        if len2 == 1 {
            if len1 <= 0 {
                return errors.New(" len1 > 0;")
            }
            dest -= len1
            cursor1 -= len1

            copy(a[dest+1:dest+1+len1], a[cursor1+1:cursor1+1+len1])
            a[dest] = tmp[cursor2] // Move first elt of run2 to front of merge
        } else if len2 == 0 {
            return errors.New("comparison method violates its general contract")
        } else {
            if len1 != 0 {
                return errors.New("len1 == 0;")
            }

            if len2 <= 0 {
                return errors.New(" len2 > 0;")
            }

            copy(a[dest-(len2-1):dest+1], tmp)
        }
        return
    }

    func (h *${timSortHandler}) ensureCapacity(minCapacity int) []${sliceType} {
        if len(h.tmp) < minCapacity {
            // Compute smallest power of 2 > minCapacity
            newSize := minCapacity
            newSize |= newSize >> 1
            newSize |= newSize >> 2
            newSize |= newSize >> 4
            newSize |= newSize >> 8
            newSize |= newSize >> 16
            newSize++

            if newSize < 0 { // Not bloody likely!
                newSize = minCapacity
            } else {
                ns := len(h.a) / 2
                if ns < newSize {
                    newSize = ns
                }
            }

            h.tmp = make([]${sliceType}, newSize)
        }
        return h.tmp
    }
    `;
}
