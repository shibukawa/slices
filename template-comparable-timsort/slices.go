package template_comparable_timsort

import (
	"errors"
	"fmt"
	"github.com/cheekybits/genny/generic"
)

type ValueType generic.Number

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

const (
	/**
	 * This is the minimum sized sequence that will be merged.  Shorter
	 * sequences will be lengthened by calling binarySort.  If the entire
	 * array is less than this length, no merges will be performed.
	 *
	 * This constant should be a power of two.  It was 64 in Tim Peter's C
	 * implementation, but 32 was empirically determined to work better in
	 * this implementation.  In the unlikely event that you set this constant
	 * to be a number that's not a power of two, you'll need to change the
	 * {@link #minRunLength} computation.
	 *
	 * If you decrease this constant, you must change the stackLen
	 * computation in the TimSort constructor, or you risk an
	 * ArrayOutOfBounds exception.  See listsort.txt for a discussion
	 * of the minimum stack length required as a function of the length
	 * of the array being sorted and the minimum merge sequence length.
	 */
	minMerge = 32
	// mk: tried higher MIN_MERGE and got slower sorting (348->375)
	//	c_MIN_MERGE = 64

	/**
	 * When we get into galloping mode, we stay there until both runs win less
	 * often than cminGallop consecutive times.
	 */
	minGallop = 7

	/**
	 * Maximum initial size of tmp array, which is used for merging.  The array
	 * can grow to accommodate demand.
	 *
	 * Unlike Tim's original C version, we do not allocate this much storage
	 * when sorting smaller arrays.  This change was required for performance.
	 */
	initialTmpStorageLength = 256
)

type timSortHandler struct {

	/**
	 * The array being sorted.
	 */
	a []ValueType

	/**
	 * This controls when we get *into* galloping mode.  It is initialized
	 * to cminGallop.  The mergeLo and mergeHi methods nudge it higher for
	 * random data, and lower for highly structured data.
	 */
	minGallop int

	/**
	 * Temp storage for merges.
	 */
	tmp []ValueType // Actual runtime type will be Object[], regardless of ValueType

	/**
	 * A stack of pending runs yet to be merged.  Run i starts at
	 * address base[i] and extends for len[i] elements.  It's always
	 * true (so long as the indices are in bounds) that:
	 *
	 *     runBase[i] + runLen[i] == runBase[i + 1]
	 *
	 * so we could cut the storage for this, but it's a minor amount,
	 * and keeping all the info explicit simplifies the code.
	 */
	stackSize int // Number of pending runs on stack
	runBase   []int
	runLen    []int
}

/**
 * Creates a TimSort instance to maintain the state of an ongoing sort.
 *
 * @param a the array to be sorted
 */
func newTimSort(a []ValueType) (h *timSortHandler) {
	h = new(timSortHandler)

	h.a = a
	h.minGallop = minGallop
	h.stackSize = 0

	// Allocate temp storage (which may be increased later if necessary)
	len := len(a)

	tmpSize := initialTmpStorageLength
	if len < 2*tmpSize {
		tmpSize = len / 2
	}

	h.tmp = make([]ValueType, tmpSize)

	/*
	 * Allocate runs-to-be-merged stack (which cannot be expanded).  The
	 * stack length requirements are described in listsort.txt.  The C
	 * version always uses the same stack length (85), but this was
	 * measured to be too expensive when sorting "mid-sized" arrays (e.g.,
	 * 100 elements) in Java.  Therefore, we use smaller (but sufficiently
	 * large) stack lengths for smaller arrays.  The "magic numbers" in the
	 * computation below must be changed if c_MIN_MERGE is decreased.  See
	 * the c_MIN_MERGE declaration above for more information.
	 */
	// mk: confirmed that for small sorts this optimization gives measurable (albeit small)
	// performance enhancement
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

// ValueTypeSort sorts an array using the provided comparator
func ValueTypeSort(a []ValueType) (err error) {
	lo := 0
	hi := len(a)
	nRemaining := hi

	if nRemaining < 2 {
		return // Arrays of size 0 and 1 are always sorted
	}

	// If array is small, do a "mini-TimSort" with no merges
	if nRemaining < minMerge {
		initRunLen, err := countRunAndMakeAscending(a, lo, hi)
		if err != nil {
			return err
		}

		return binarySort(a, lo, hi, lo+initRunLen)
	}

	/**
	 * March over the array once, left to right, finding natural runs,
	 * extending short natural runs to minRun elements, and merging runs
	 * to maintain stack invariant.
	 */

	ts := newTimSort(a)
	minRun, err := minRunLength(nRemaining)
	if err != nil {
		return
	}
	for {
		// Identify next run
		runLen, err := countRunAndMakeAscending(a, lo, hi)
		if err != nil {
			return err
		}

		// If run is short, extend to min(minRun, nRemaining)
		if runLen < minRun {
			force := minRun
			if nRemaining <= minRun {
				force = nRemaining
			}
			if err = binarySort(a, lo, lo+force, lo+runLen); err != nil {
				return err
			}
			runLen = force
		}

		// Push run onto pending-run stack, and maybe merge
		ts.pushRun(lo, runLen)
		if err = ts.mergeCollapse(); err != nil {
			return err
		}

		// Advance to find next run
		lo += runLen
		nRemaining -= runLen
		if nRemaining == 0 {
			break
		}
	}

	// Merge all remaining runs to complete sort
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

/**
 * Sorts the specified portion of the specified array using a binary
 * insertion sort.  This is the best method for sorting small numbers
 * of elements.  It requires O(n log n) compares, but O(n^2) data
 * movement (worst case).
 *
 * If the initial part of the specified range is already sorted,
 * this method can take advantage of it: the method assumes that the
 * elements from index {@code lo}, inclusive, to {@code start},
 * exclusive are already sorted.
 *
 * @param a the array in which a range is to be sorted
 * @param lo the index of the first element in the range to be sorted
 * @param hi the index after the last element in the range to be sorted
 * @param start the index of the first element in the range that is
 *        not already known to be sorted (@code lo <= start <= hi}
 * @param c comparator to used for the sort
 */
func binarySort(a []ValueType, lo, hi, start int) (err error) {
	if lo > start || start > hi {
		return errors.New("lo <= start && start <= hi")
	}

	if start == lo {
		start++
	}

	for ; start < hi; start++ {
		pivot := a[start]

		// Set left (and right) to the index where a[start] (pivot) belongs
		left := lo
		right := start

		if left > right {
			return errors.New("left <= right")
		}

		/*
		 * Invariants:
		 *   pivot >= all in [lo, left).
		 *   pivot <  all in [right, start).
		 */
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

		/*
		 * The invariants still hold: pivot >= all in [lo, left) and
		 * pivot < all in [left, start), so pivot belongs at left.  Note
		 * that if there are elements equal to pivot, left points to the
		 * first slot after them -- that's why this sort is stable.
		 * Slide elements over to make room to make room for pivot.
		 */
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

/**
  * Returns the length of the run beginning at the specified position in
  * the specified array and reverses the run if it is descending (ensuring
  * that the run will always be ascending when the method returns).
  *
  * A run is the longest ascending sequence with:
  *
  *    a[lo] <= a[lo + 1] <= a[lo + 2] <= ...
  *
  * or the longest descending sequence with:
  *
  *    a[lo] >  a[lo + 1] >  a[lo + 2] >  ...
  *
  * For its intended use in a stable mergesort, the strictness of the
  * definition of "descending" is needed so that the call can safely
  * reverse a descending sequence without violating stability.
  *
  * @param a the array in which a run is to be counted and possibly reversed
  * @param lo index of the first element in the run
  * @param hi index after the last element that may be contained in the run.
           It is required that @code{lo < hi}.
  * @param c the comparator to used for the sort
  * @return  the length of the run beginning at the specified position in
  *          the specified array
*/
func countRunAndMakeAscending(a []ValueType, lo, hi int) (int, error) {

	if lo >= hi {
		return 0, errors.New("lo < hi")
	}

	runHi := lo + 1
	if runHi == hi {
		return 1, nil
	}

	// Find end of run, and reverse range if descending
	if a[runHi] < a[lo] { // Descending
		runHi++

		for runHi < hi && a[runHi] < a[runHi-1] {
			runHi++
		}
		reverseRange(a, lo, runHi)
	} else { // Ascending
		for runHi < hi && !(a[runHi] < a[runHi-1]) {
			runHi++
		}
	}

	return runHi - lo, nil
}

/**
 * Reverse the specified range of the specified array.
 *
 * @param a the array in which a range is to be reversed
 * @param lo the index of the first element in the range to be reversed
 * @param hi the index after the last element in the range to be reversed
 */
func reverseRange(a []ValueType, lo, hi int) {
	hi--
	for lo < hi {
		a[lo], a[hi] = a[hi], a[lo]
		lo++
		hi--
	}
}

/**
 * Returns the minimum acceptable run length for an array of the specified
 * length. Natural runs shorter than this will be extended with
 * {@link #binarySort}.
 *
 * Roughly speaking, the computation is:
 *
 *  If n < c_MIN_MERGE, return n (it's too small to bother with fancy stuff).
 *  Else if n is an exact power of 2, return c_MIN_MERGE/2.
 *  Else return an int k, c_MIN_MERGE/2 <= k <= c_MIN_MERGE, such that n/k
 *   is close to, but strictly less than, an exact power of 2.
 *
 * For the rationale, see listsort.txt.
 *
 * @param n the length of the array to be sorted
 * @return the length of the minimum run to be merged
 */
func minRunLength(n int) (int, error) {
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

/**
 * Pushes the specified run onto the pending-run stack.
 *
 * @param runBase index of the first element in the run
 * @param runLen  the number of elements in the run
 */
func (h *timSortHandler) pushRun(runBase, runLen int) {
	h.runBase[h.stackSize] = runBase
	h.runLen[h.stackSize] = runLen
	h.stackSize++
}

/**
 * Examines the stack of runs waiting to be merged and merges adjacent runs
 * until the stack invariants are reestablished:
 *
 *     1. runLen[i - 3] > runLen[i - 2] + runLen[i - 1]
 *     2. runLen[i - 2] > runLen[i - 1]
 *
 * This method is called each time a new run is pushed onto the stack,
 * so the invariants are guaranteed to hold for i < stackSize upon
 * entry to the method.
 */
func (h *timSortHandler) mergeCollapse() (err error) {
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

/**
 * Merges all runs on the stack until only one remains.  This method is
 * called once, to complete the sort.
 */
func (h *timSortHandler) mergeForceCollapse() (err error) {
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

/**
 * Merges the two runs at stack indices i and i+1.  Run i must be
 * the penultimate or antepenultimate run on the stack.  In other words,
 * i must be equal to stackSize-2 or stackSize-3.
 *
 * @param i stack index of the first of the two runs to merge
 */
func (h *timSortHandler) mergeAt(i int) (err error) {
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

	/*
	 * Record the length of the combined runs; if i is the 3rd-last
	 * run now, also slide over the last run (which isn't involved
	 * in this merge).  The current run (i+1) goes away in any case.
	 */
	h.runLen[i] = len1 + len2
	if i == h.stackSize-3 {
		h.runBase[i+1] = h.runBase[i+2]
		h.runLen[i+1] = h.runLen[i+2]
	}
	h.stackSize--

	/*
	 * Find where the first element of run2 goes in run1. Prior elements
	 * in run1 can be ignored (because they're already in place).
	 */
	k, err := gallopRight(h.a[base2], h.a, base1, len1, 0)
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

	/*
	 * Find where the last element of run1 goes in run2. Subsequent elements
	 * in run2 can be ignored (because they're already in place).
	 */
	len2, err = gallopLeft(h.a[base1+len1-1], h.a, base2, len2, len2-1)
	if err != nil {
		return
	}
	if len2 < 0 {
		return errors.New(" len2 >= 0;")
	}
	if len2 == 0 {
		return
	}

	// Merge remaining runs, using tmp array with min(len1, len2) elements
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

/**
 * Locates the position at which to insert the specified key into the
 * specified sorted range; if the range contains an element equal to key,
 * returns the index of the leftmost equal element.
 *
 * @param key the key whose insertion point to search for
 * @param a the array in which to search
 * @param base the index of the first element in the range
 * @param len the length of the range; must be > 0
 * @param hint the index at which to begin the search, 0 <= hint < n.
 *     The closer hint is to the result, the faster this method will run.
 * @return the int k,  0 <= k <= n such that a[b + k - 1] < key <= a[b + k],
 *    pretending that a[b - 1] is minus infinity and a[b + n] is infinity.
 *    In other words, key belongs at index b + k; or in other words,
 *    the first k elements of a should precede key, and the last n - k
 *    should follow it.
 */
func gallopLeft(key ValueType, a []ValueType, base, len, hint int) (int, error) {
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

		// Make offsets relative to base
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

		// Make offsets relative to base
		tmp := lastOfs
		lastOfs = hint - ofs
		ofs = hint - tmp
	}

	if -1 > lastOfs || lastOfs >= ofs || ofs > len {
		return 0, errors.New(" -1 <= lastOfs && lastOfs < ofs && ofs <= len;")
	}

	/*
	 * Now a[base+lastOfs] < key <= a[base+ofs], so key belongs somewhere
	 * to the right of lastOfs but no farther right than ofs.  Do a binary
	 * search, with invariant a[base + lastOfs - 1] < key <= a[base + ofs].
	 */
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

/**
 * Like gallopLeft, except that if the range contains an element equal to
 * key, gallopRight returns the index after the rightmost equal element.
 *
 * @param key the key whose insertion point to search for
 * @param a the array in which to search
 * @param base the index of the first element in the range
 * @param len the length of the range; must be > 0
 * @param hint the index at which to begin the search, 0 <= hint < n.
 *     The closer hint is to the result, the faster this method will run.
 * @return the int k,  0 <= k <= n such that a[b + k - 1] <= key < a[b + k]
 */
func gallopRight(key ValueType, a []ValueType, base, len, hint int) (int, error) {
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

		// Make offsets relative to b
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

		// Make offsets relative to b
		lastOfs += hint
		ofs += hint
	}
	if -1 > lastOfs || lastOfs >= ofs || ofs > len {
		return 0, errors.New("-1 <= lastOfs && lastOfs < ofs && ofs <= len")
	}

	/*
	 * Now a[b + lastOfs] <= key < a[b + ofs], so key belongs somewhere to
	 * the right of lastOfs but no farther right than ofs.  Do a binary
	 * search, with invariant a[b + lastOfs - 1] <= key < a[b + ofs].
	 */
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

/**
 * Merges two adjacent runs in place, in a stable fashion.  The first
 * element of the first run must be greater than the first element of the
 * second run (a[base1] > a[base2]), and the last element of the first run
 * (a[base1 + len1-1]) must be greater than all elements of the second run.
 *
 * For performance, this method should be called only when len1 <= len2;
 * its twin, mergeHi should be called if len1 >= len2.  (Either method
 * may be called if len1 == len2.)
 *
 * @param base1 index of first element in first run to be merged
 * @param len1  length of first run to be merged (must be > 0)
 * @param base2 index of first element in second run to be merged
 *        (must be aBase + aLen)
 * @param len2  length of second run to be merged (must be > 0)
 */
func (h *timSortHandler) mergeLo(base1, len1, base2, len2 int) (err error) {
	if len1 <= 0 || len2 <= 0 || base1+len1 != base2 {
		return errors.New(" len1 > 0 && len2 > 0 && base1 + len1 == base2")
	}

	// Copy first run into temp array
	a := h.a // For performance
	tmp := h.ensureCapacity(len1)

	copy(tmp, a[base1:base1+len1])

	cursor1 := 0     // Indexes into tmp array
	cursor2 := base2 // Indexes int a
	dest := base1    // Indexes int a

	// Move first element of second run and deal with degenerate cases
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

	minGallop := h.minGallop //  "    "       "     "      "

outer:
	for {
		count1 := 0 // Number of times in a row that first run won
		count2 := 0 // Number of times in a row that second run won

		/*
		 * Do the straightforward thing until (if ever) one run starts
		 * winning consistently.
		 */
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

		/*
		 * One run is winning so consistently that galloping may be a
		 * huge win. So try that, and continue galloping until (if ever)
		 * neither run appears to be winning consistently anymore.
		 */
		for {
			if len1 <= 1 || len2 <= 0 {
				return errors.New("len1 > 1 && len2 > 0")
			}
			count1, err = gallopRight(a[cursor2], tmp, cursor1, len1, 0)
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

			count2, err = gallopLeft(tmp[cursor1], a, cursor2, len2, 0)
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
	} // End of "outer" loop

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

/**
 * Like mergeLo, except that this method should be called only if
 * len1 >= len2; mergeLo should be called if len1 <= len2.  (Either method
 * may be called if len1 == len2.)
 *
 * @param base1 index of first element in first run to be merged
 * @param len1  length of first run to be merged (must be > 0)
 * @param base2 index of first element in second run to be merged
 *        (must be aBase + aLen)
 * @param len2  length of second run to be merged (must be > 0)
 */
func (h *timSortHandler) mergeHi(base1, len1, base2, len2 int) (err error) {
	if len1 <= 0 || len2 <= 0 || base1+len1 != base2 {
		return errors.New("len1 > 0 && len2 > 0 && base1 + len1 == base2;")
	}

	// Copy second run into temp array
	a := h.a // For performance
	tmp := h.ensureCapacity(len2)

	copy(tmp, a[base2:base2+len2])

	cursor1 := base1 + len1 - 1 // Indexes into a
	cursor2 := len2 - 1         // Indexes into tmp array
	dest := base2 + len2 - 1    // Indexes into a

	// Move last element of first run and deal with degenerate cases
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

	minGallop := h.minGallop //  "    "       "     "      "

outer:
	for {
		count1 := 0 // Number of times in a row that first run won
		count2 := 0 // Number of times in a row that second run won

		/*
		 * Do the straightforward thing until (if ever) one run
		 * appears to win consistently.
		 */
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

		/*
		 * One run is winning so consistently that galloping may be a
		 * huge win. So try that, and continue galloping until (if ever)
		 * neither run appears to be winning consistently anymore.
		 */
		for {
			if len1 <= 0 || len2 <= 1 {
				return errors.New(" len1 > 0 && len2 > 1;")
			}
			if gr, err := gallopRight(tmp[cursor2], a, base1, len1, len1-1); err == nil {
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

			if gl, err := gallopLeft(a[cursor1], tmp, 0, len2, len2-1); err == nil {
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

/**
 * Ensures that the external array tmp has at least the specified
 * number of elements, increasing its size if necessary.  The size
 * increases exponentially to ensure amortized linear time complexity.
 *
 * @param minCapacity the minimum required capacity of the tmp array
 * @return tmp, whether or not it grew
 */
func (h *timSortHandler) ensureCapacity(minCapacity int) []ValueType {
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

		h.tmp = make([]ValueType, newSize)
	}

	return h.tmp
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
