package autoimport

import (
	"maps"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/microsoft/typescript-go/internal/core"
)

type Trie[T any] struct {
	root *trieNode[T]
}

func (t *Trie[T]) Search(s string) []*T {
	s = strings.ToLower(s)
	if t.root == nil {
		return nil
	}
	node := t.root
	for _, r := range s {
		if node.children[r] == nil {
			return nil
		}
		node = node.children[r]
	}

	var results []*T
	results = append(results, node.values...)
	for _, child := range node.children {
		results = append(results, child.collectValues()...)
	}
	return results
}

type trieNode[T any] struct {
	children map[rune]*trieNode[T]
	values   []*T
}

func (n *trieNode[T]) clone() *trieNode[T] {
	newNode := &trieNode[T]{
		children: maps.Clone(n.children),
		values:   slices.Clone(n.values),
	}
	return newNode
}

func (n *trieNode[T]) collectValues() []*T {
	var results []*T
	results = append(results, n.values...)
	for _, child := range n.children {
		results = append(results, child.collectValues()...)
	}
	return results
}

type TrieBuilder[T any] struct {
	t      *Trie[T]
	cloned map[*trieNode[T]]struct{}
}

func NewTrieBuilder[T any](trie *Trie[T]) *TrieBuilder[T] {
	return &TrieBuilder[T]{
		t:      trie,
		cloned: make(map[*trieNode[T]]struct{}),
	}
}

func (t *TrieBuilder[T]) cloneNode(n *trieNode[T]) *trieNode[T] {
	if _, ok := t.cloned[n]; ok {
		return n
	}
	clone := n.clone()
	t.cloned[n] = struct{}{}
	return clone
}

func (t *TrieBuilder[T]) Trie() *Trie[T] {
	trie := t.t
	t.t = nil
	return trie
}

func (t *TrieBuilder[T]) Insert(s string, value *T) {
	if t.t == nil {
		panic("insert called after TrieBuilder.Trie()")
	}
	if t.t.root == nil {
		t.t.root = &trieNode[T]{children: make(map[rune]*trieNode[T])}
	}

	node := t.t.root
	for _, r := range s {
		r = unicode.ToLower(r)
		if node.children[r] == nil {
			child := &trieNode[T]{children: make(map[rune]*trieNode[T])}
			node.children[r] = child
			t.cloned[child] = struct{}{}
			node = child
		} else {
			node = t.cloneNode(node.children[r])
		}
	}
	node.values = append(node.values, value)
}

func (t *TrieBuilder[T]) InsertAsWords(s string, value *T) {
	indices := wordIndices(s)
	for _, start := range indices {
		t.Insert(s[start:], value)
	}
}

// wordIndices splits an identifier into its constituent words based on camelCase and snake_case conventions
// by returning the starting byte indices of each word.
//   - CamelCase
//     ^    ^
//   - snake_case
//     ^     ^
//   - ParseURL
//     ^    ^
//   - __proto__
//     ^
func wordIndices(s string) []int {
	var indices []int
	for byteIndex, runeValue := range s {
		if byteIndex == 0 {
			indices = append(indices, byteIndex)
			continue
		}
		if runeValue == '_' {
			if byteIndex+1 < len(s) && s[byteIndex+1] != '_' {
				indices = append(indices, byteIndex+1)
			}
			continue
		}
		if isUpper(runeValue) && isLower(core.FirstResult(utf8.DecodeLastRuneInString(s[:byteIndex]))) || (byteIndex+1 < len(s) && isLower(core.FirstResult(utf8.DecodeRuneInString(s[byteIndex+1:])))) {
			indices = append(indices, byteIndex)
		}
	}
	return indices
}

func isUpper(c rune) bool {
	return c >= 'A' && c <= 'Z'
}

func isLower(c rune) bool {
	return c >= 'a' && c <= 'z'
}
