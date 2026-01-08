package core

import "strings"

type TextChange struct {
	TextRange
	NewText string
}

func (t TextChange) ApplyTo(text string) string {
	return text[:t.Pos()] + t.NewText + text[t.End():]
}

func ApplyBulkEdits(text string, edits []TextChange) string {
	b := strings.Builder{}
	b.Grow(len(text))
	lastEnd := 0
	for _, e := range edits {
		start := e.Pos()
		if start != lastEnd {
			b.WriteString(text[lastEnd:e.Pos()])
		}
		b.WriteString(e.NewText)

		lastEnd = e.End()
	}
	b.WriteString(text[lastEnd:])

	return b.String()
}
