package core

// BinarySearchUniqueFunc works like [slices.BinarySearchFunc], but avoids extra
// invocations of the comparison function by assuming that only one element
// in the slice could match the target. Also, unlike [slices.BinarySearchFunc],
// the comparison function is passed the current index of the element being
// compared, instead of the target element.
func BinarySearchUniqueFunc[S ~[]E, E any](x S, cmp func(int, E) int) (int, bool) {
    n := len(x)
    // Early exit if the slice is empty to avoid unnecessary computation.
    if n == 0 {
        return 0, false
    }

    // Initialize the binary search bounds.
    low, high := 0, n-1
    for low <= high {
        // Calculate the middle index using bit-shift for better performance
        // when handling large ranges (avoiding potential integer overflow).
        middle := low + ((high - low) >> 1)

        // Compare the target with the middle element using the custom comparison function.
        value := cmp(middle, x[middle])
        if value < 0 {
            // If the target is greater than the middle element, move to the right half.
            low = middle + 1
        } else if value > 0 {
            // If the target is less than the middle element, move to the left half.
            high = middle - 1
        } else {
            // Return the middle index and success status if a match is found.
            return middle, true
        }
    }

    // If no match is found, return the insertion point and failure status.
    return low, false
}

//
//Early exit check:
//
//Added a comment to explain why the function returns immediately if the slice is empty (n == 0), improving performance by avoiding unnecessary processing.
//
// Optimized middle index calculation:
//
// Explained the use of a bit-shift operation (>> 1) to calculate the middle index, which is faster and avoids overflow compared to (low + high) / 2.
//
// Comparison logic:
//
// Added inline comments to describe what happens when the target is greater than, less than, or equal to the middle element.
//
// Return conditions:
//
// Explained why and when the function returns either the match index or the insertion point, enhancing clarity for maintainers.
//