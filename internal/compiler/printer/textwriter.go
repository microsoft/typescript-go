package printer

import (
	"strings"
	"unicode/utf8"

	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/compiler/stringutil"
)

var _ EmitTextWriter = &textWriter{}

type textWriter struct {
	newLine                 string
	builder                 strings.Builder
	lastWritten             string
	indent                  int
	lineStart               bool
	lineCount               int
	linePos                 int
	hasTrailingCommentState bool
}

// clear implements EmitTextWriter.
func (w *textWriter) clear() {
	// Is it worth reusing the old string builder?
	w.builder.Reset()
	w.lastWritten = ""
	w.indent = 0
	w.lineStart = true
	w.lineCount = 0
	w.linePos = 0
	w.hasTrailingCommentState = false
}

// decreaseIndent implements EmitTextWriter.
func (w *textWriter) decreaseIndent() {
	w.indent--
}

// getColumn implements EmitTextWriter.
func (w *textWriter) getColumn() int {
	if w.lineStart {
		return w.indent * 4
	}
	return w.builder.Len() - w.linePos
}

// getIndent implements EmitTextWriter.
func (w *textWriter) getIndent() int {
	return w.indent
}

// getLine implements EmitTextWriter.
func (w *textWriter) getLine() int {
	return w.lineCount
}

// getText implements EmitTextWriter.
func (w *textWriter) getText() string {
	return w.builder.String()
}

// getTextPos implements EmitTextWriter.
func (w *textWriter) getTextPos() int {
	return w.builder.Len()
}

// hasTrailingComment implements EmitTextWriter.
func (w textWriter) hasTrailingComment() bool {
	return w.hasTrailingCommentState
}

// hasTrailingWhitespace implements EmitTextWriter.
func (w *textWriter) hasTrailingWhitespace() bool {
	if w.builder.Len() == 0 {
		return false
	}
	ch, _ := utf8.DecodeLastRuneInString(w.lastWritten)
	if ch == utf8.RuneError {
		return false
	}
	return stringutil.IsWhiteSpaceLike(ch)
}

// increaseIndent implements EmitTextWriter.
func (w *textWriter) increaseIndent() {
	w.indent++
}

// isAtStartOfLine implements EmitTextWriter.
func (w *textWriter) isAtStartOfLine() bool {
	return w.lineStart
}

// rawWrite implements EmitTextWriter.
func (w *textWriter) rawWrite(s string) {
	if s != "" {
		w.builder.WriteString(s)
		w.lastWritten = s
		w.updateLineCountAndPosFor(s)
		w.hasTrailingCommentState = false
	}
}

func (w *textWriter) updateLineCountAndPosFor(s string) {
	lineStartsOfS := stringutil.ComputeLineStarts(s)
	if len(lineStartsOfS) > 1 {
		w.lineCount += len(lineStartsOfS) - 1
		curLen := w.builder.Len()
		w.linePos = curLen - len(s) + int(lineStartsOfS[len(lineStartsOfS)-1])
		w.lineStart = (w.linePos - curLen) == 0
		return
	}
	w.lineStart = false
}

func getIndentString(indent int) string {
	switch indent {
	case 0:
		return ""
	case 1:
		return "    "
	default:
		// TODO: This is cached in tsc - should it be cached here?
		return strings.Repeat("    ", indent)
	}
}

func (w *textWriter) writeText(s string) {
	if s != "" {
		if w.lineStart {
			w.builder.WriteString(getIndentString(w.indent))
			w.lineStart = false
		}
		w.builder.WriteString(s)
		w.lastWritten = s
		w.updateLineCountAndPosFor(s)
	}
}

// write implements EmitTextWriter.
func (w *textWriter) write(s string) {
	if s != "" {
		w.hasTrailingCommentState = false
	}
	w.writeText(s)
}

// writeComment implements EmitTextWriter.
func (w *textWriter) writeComment(text string) {
	if text != "" {
		w.hasTrailingCommentState = true
	}
	w.writeText(text)
}

// writeKeyword implements EmitTextWriter.
func (w *textWriter) writeKeyword(text string) {
	w.write(text)
}

// writeLine implements EmitTextWriter.
func (w *textWriter) writeLine(force ...bool) {
	if !w.lineStart || force[0] == true {
		w.builder.WriteString(w.newLine)
		w.lastWritten = w.newLine
		w.lineCount++
		w.linePos = w.builder.Len()
		w.lineStart = true
		w.hasTrailingCommentState = false
	}
}

// writeLiteral implements EmitTextWriter.
func (w *textWriter) writeLiteral(s string) {
	w.write(s)
}

// writeOperator implements EmitTextWriter.
func (w *textWriter) writeOperator(text string) {
	w.write(text)
}

// writeParameter implements EmitTextWriter.
func (w *textWriter) writeParameter(text string) {
	w.write(text)
}

// writeProperty implements EmitTextWriter.
func (w *textWriter) writeProperty(text string) {
	w.write(text)
}

// writePunctuation implements EmitTextWriter.
func (w *textWriter) writePunctuation(text string) {
	w.write(text)
}

// writeSpace implements EmitTextWriter.
func (w *textWriter) writeSpace(text string) {
	w.write(text)
}

// writeStringLiteral implements EmitTextWriter.
func (w *textWriter) writeStringLiteral(text string) {
	w.write(text)
}

// writeSymbol implements EmitTextWriter.
func (w *textWriter) writeSymbol(text string, symbol compiler.Symbol) {
	w.write(text)
}

// writeTrailingSemicolon implements EmitTextWriter.
func (w *textWriter) writeTrailingSemicolon(text string) {
	w.write(text)
}

func NewTextWriter(newLine string) EmitTextWriter {
	return &textWriter{newLine: newLine}
}
