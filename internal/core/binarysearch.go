package core

// BinarySearchUniqueFunc works like [slices.BinarySearchFunc], but avoids extra
// invocations of the comparison function by assuming that only one element
// in the slice could match the target. Also, unlike [slices.BinarySearchFunc],
// the comparison function is passed the current index of the element being
// compared, instead of the target element.
func BinarySearchUniqueFunc[S ~[]E, E, T any](x S, target T, cmp func(int, E) int) (int, bool) {
	n := len(x)
	low, high := 0, n-1
	for low <= high {
		middle := low + ((high - low) >> 1)
		switch cmp(middle, x[middle]) {
		case -1:
			low = middle + 1
		case 0:
			return middle, true
		case 1:
			high = middle - 1
		}
	}
	return low, false
}
