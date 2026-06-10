package stringutil

func isInRuneRanges(cp rune, ranges []rune) bool {
	// Bail out quickly if it couldn't possibly be in the map
	if cp < ranges[0] {
		return false
	}
	// Perform binary search in one of the flattened range maps.
	lo := 0
	hi := len(ranges)
	for lo+1 < hi {
		mid := lo + (hi-lo)/2
		// mid has to be even to catch the beginning of a range.
		mid -= mid % 2
		if ranges[mid] <= cp && cp <= ranges[mid+1] {
			return true
		}
		if cp < ranges[mid] {
			hi = mid
		} else {
			lo = mid + 2
		}
	}
	return false
}
