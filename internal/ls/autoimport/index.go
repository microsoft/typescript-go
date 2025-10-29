package autoimport

import (
	"slices"
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

	prefix = strings.ToLower(prefix)
	if len(prefix) == 0 {
		return nil
	}

	// Get the first rune of the prefix
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

// IndexBuilder builds an Index with copy-on-write semantics for efficient updates.
type IndexBuilder[T Named] struct {
	idx    *Index[T]
	cloned bool
}

// NewIndexBuilder creates a new IndexBuilder from an existing Index.
// If idx is nil, a new empty Index will be created.
func NewIndexBuilder[T Named](idx *Index[T]) *IndexBuilder[T] {
	if idx == nil {
		idx = &Index[T]{
			entries: make([]T, 0),
			index:   make(map[rune][]int),
		}
	}
	return &IndexBuilder[T]{
		idx:    idx,
		cloned: false,
	}
}

func (b *IndexBuilder[T]) ensureCloned() {
	if !b.cloned {
		newIdx := &Index[T]{
			entries: slices.Clone(b.idx.entries),
			index:   make(map[rune][]int, len(b.idx.index)),
		}
		for k, v := range b.idx.index {
			newIdx.index[k] = slices.Clone(v)
		}
		b.idx = newIdx
		b.cloned = true
	}
}

// Insert adds a value to the index.
// The value will be indexed by the first letter of its name.
func (b *IndexBuilder[T]) Insert(value T) {
	if b.idx == nil {
		panic("insert called after IndexBuilder.Index()")
	}
	b.ensureCloned()

	name := value.Name()
	name = strings.ToLower(name)
	if len(name) == 0 {
		return
	}

	firstRune, _ := utf8.DecodeRuneInString(name)
	if firstRune == utf8.RuneError {
		return
	}

	entryIndex := len(b.idx.entries)
	b.idx.entries = append(b.idx.entries, value)
	b.idx.index[firstRune] = append(b.idx.index[firstRune], entryIndex)
}

// InsertAsWords adds a value to the index, indexing it by the first letter of each word
// in its name. Words are determined by camelCase, PascalCase, and snake_case conventions.
func (b *IndexBuilder[T]) InsertAsWords(value T) {
	if b.idx == nil {
		panic("insert called after IndexBuilder.Index()")
	}
	b.ensureCloned()

	name := value.Name()
	entryIndex := len(b.idx.entries)
	b.idx.entries = append(b.idx.entries, value)

	// Get all word start positions
	indices := wordIndices(name)
	seenRunes := make(map[rune]bool)

	for _, start := range indices {
		substr := name[start:]
		firstRune, _ := utf8.DecodeRuneInString(substr)
		if firstRune == utf8.RuneError {
			continue
		}
		firstRune = unicode.ToLower(firstRune)

		// Only add each letter once per entry
		if !seenRunes[firstRune] {
			b.idx.index[firstRune] = append(b.idx.index[firstRune], entryIndex)
			seenRunes[firstRune] = true
		}
	}
}

// Index returns the built Index and invalidates the builder.
func (b *IndexBuilder[T]) Index() *Index[T] {
	idx := b.idx
	b.idx = nil
	return idx
}
