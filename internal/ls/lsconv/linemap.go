package lsconv

import (
	"cmp"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/microsoft/typescript-go/internal/core"
)

type LSPLineStarts []core.TextPos

type LSPLineMap struct {
	LineStarts       LSPLineStarts
	nonAsciiLineNums []int // Sorted slice of line numbers (0-indexed) that contain non-ASCII characters
}

func ComputeLSPLineStarts(text string) *LSPLineMap {
	// This is like core.ComputeLineStarts, but only considers "\n", "\r", and "\r\n" as line breaks,
	// and records which lines contain non-ASCII characters.
	lineStarts := make([]core.TextPos, 0, strings.Count(text, "\n")+1)
	var nonAsciiLineNums []int

	textLen := core.TextPos(len(text))
	var pos core.TextPos
	var lineStart core.TextPos
	currentLineNum := 0
	currentLineHasNonAscii := false

	for pos < textLen {
		b := text[pos]
		if b < utf8.RuneSelf {
			pos++
			switch b {
			case '\r':
				if pos < textLen && text[pos] == '\n' {
					pos++
				}
				fallthrough
			case '\n':
				lineStarts = append(lineStarts, lineStart)
				if currentLineHasNonAscii {
					nonAsciiLineNums = append(nonAsciiLineNums, currentLineNum)
				}
				lineStart = pos
				currentLineNum++
				currentLineHasNonAscii = false
			}
		} else {
			_, size := utf8.DecodeRuneInString(text[pos:])
			pos += core.TextPos(size)
			currentLineHasNonAscii = true
		}
	}
	lineStarts = append(lineStarts, lineStart)
	if currentLineHasNonAscii {
		nonAsciiLineNums = append(nonAsciiLineNums, currentLineNum)
	}

	return &LSPLineMap{
		LineStarts:       lineStarts,
		nonAsciiLineNums: nonAsciiLineNums,
	}
}

func (lm *LSPLineMap) ComputeIndexOfLineStart(targetPos core.TextPos) int {
	// port of computeLineOfPosition(lineStarts: readonly number[], position: number, lowerBound?: number): number {
	lineNumber, ok := slices.BinarySearchFunc(lm.LineStarts, targetPos, func(p, t core.TextPos) int {
		return cmp.Compare(int(p), int(t))
	})
	if !ok && lineNumber > 0 {
		// If the actual position was not found, the binary search returns where the target line start would be inserted
		// if the target was in the slice.
		// e.g. if the line starts at [5, 10, 23, 80] and the position requested was 20
		// then the search will return (3, false).
		//
		// We want the index of the previous line start, so we subtract 1.
		lineNumber = lineNumber - 1
	}
	return lineNumber
}

func (lm *LSPLineMap) IsLineAsciiOnly(lineNum int) bool {
	_, found := slices.BinarySearch(lm.nonAsciiLineNums, lineNum)
	return !found
}
