package autoimport

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// Named is a constraint for types that can provide their name.
type Named interface {
	Name() string
}

// Index stores entries with an index mapping lowercase letters to entries whose name
// has a word starting with that letter. This supports efficient fuzzy matching.
type Index[T Named] struct {
	entries []T
	index   map[rune][]int
}

// Search returns all entries whose name contains the characters of prefix in order.
// The search first uses the index to narrow down candidates by the first letter,
// then filters by checking if the name contains all characters in order.
func (idx *Index[T]) Search(prefix string) []T {
	if idx == nil || len(idx.entries) == 0 {
		return nil
	}

	if len(prefix) == 0 {
		return idx.entries
	}

	prefix = strings.ToLower(prefix)
	firstRune, _ := utf8.DecodeRuneInString(prefix)
	if firstRune == utf8.RuneError {
		return nil
	}
	firstRune = unicode.ToLower(firstRune)

	// Look up entries that have words starting with this letter
	indices, ok := idx.index[firstRune]
	if !ok {
		return nil
	}

	// Filter entries by checking if they contain all characters in order
	results := make([]T, 0, len(indices))
	for _, i := range indices {
		entry := idx.entries[i]
		if containsCharsInOrder(entry.Name(), prefix) {
			results = append(results, entry)
		}
	}
	return results
}

// containsCharsInOrder checks if str contains all characters from pattern in order (case-insensitive).
func containsCharsInOrder(str, pattern string) bool {
	str = strings.ToLower(str)
	pattern = strings.ToLower(pattern)

	patternIdx := 0
	for _, ch := range str {
		if patternIdx < len(pattern) {
			patternRune, size := utf8.DecodeRuneInString(pattern[patternIdx:])
			if ch == patternRune {
				patternIdx += size
			}
		}
	}
	return patternIdx == len(pattern)
}

// insertAsWords adds a value to the index keyed by the first letter of each word in its name.
func (idx *Index[T]) insertAsWords(value T) {
	if idx.index == nil {
		idx.index = make(map[rune][]int)
	}

	name := value.Name()
	entryIndex := len(idx.entries)
	idx.entries = append(idx.entries, value)

	indices := wordIndices(name)
	seenRunes := make(map[rune]bool)

	for _, start := range indices {
		substr := name[start:]
		firstRune, _ := utf8.DecodeRuneInString(substr)
		if firstRune == utf8.RuneError {
			continue
		}
		firstRune = unicode.ToLower(firstRune)

		if !seenRunes[firstRune] {
			idx.index[firstRune] = append(idx.index[firstRune], entryIndex)
			seenRunes[firstRune] = true
		}
	}
}
