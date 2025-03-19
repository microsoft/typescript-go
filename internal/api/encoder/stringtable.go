package encoder

import (
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
)

type stringTable struct {
	fileText     string
	otherStrings *strings.Builder
	// offsets are pos/end pairs
	offsets []uint32
}

func newStringTable(fileText string, stringCount int) *stringTable {
	builder := &strings.Builder{}
	return &stringTable{
		fileText:     fileText,
		otherStrings: builder,
		offsets:      make([]uint32, 0, stringCount*2),
	}
}

func (t *stringTable) add(text string, kind ast.Kind, pos int, end int) uint32 {
	length := len(text)
	// pos includes leading trivia, but we can usually infer the actual start of the
	// string from the kind and end
	endOffset := 0
	if kind == ast.KindStringLiteral || kind == ast.KindTemplateTail || kind == ast.KindNoSubstitutionTemplateLiteral {
		endOffset = 1
	}
	end = end - endOffset
	start := end - length
	fileSlice := t.fileText[start:end]
	if fileSlice == text {
		index := len(t.offsets)
		t.offsets = append(t.offsets, uint32(start), uint32(end))
		return uint32(index)
	}
	// no exact match, so we need to add it to the string table
	offset := len(t.fileText) + t.otherStrings.Len()
	t.otherStrings.WriteString(text)
	index := len(t.offsets)
	t.offsets = append(t.offsets, uint32(offset), uint32(offset+length))
	return uint32(index)
}

func (t *stringTable) encode() (result []byte, err error) {
	result = make([]byte, 0, t.encodedLength())
	if result, err = appendUint32s(result, t.offsets...); err != nil {
		return nil, err
	}
	result = append(result, t.fileText...)
	result = append(result, t.otherStrings.String()...)
	return result, nil
}

func (t *stringTable) stringLength() int {
	return len(t.fileText) + t.otherStrings.Len()
}

func (t *stringTable) encodedLength() int {
	return len(t.offsets)*4 + len(t.fileText) + t.otherStrings.Len()
}
