package lsptestutil

import (
	"strings"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls/lsconv"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
)

type LsScript struct {
	file string
	text string
}

func NewLsScript(file string, text string) *LsScript {
	return &LsScript{file: file, text: text}
}

var _ lsconv.Script = (*LsScript)(nil)

func (s *LsScript) FileName() string { return s.file }
func (s *LsScript) Text() string     { return s.text }

func PositionToLineAndCharacter(file string, text string, substring string, index int) lsproto.Position {
	offset := nthIndexOf(text, substring, index)

	converters := lsconv.NewConverters(lsproto.PositionEncodingKindUTF8, func(fileName string) *lsconv.LSPLineMap {
		return lsconv.ComputeLSPLineStarts(text)
	})
	return converters.PositionToLineAndCharacter(NewLsScript(file, text), core.TextPos(offset))
}

func nthIndexOf(str string, substr string, n int) int {
	index := 0
	for i := range n + 1 {
		start := core.IfElse(i == 0, index, index+len(substr))
		index = strings.Index(str[start:], substr)
		if index == -1 {
			return -1
		}
		index += start
	}
	return index
}
