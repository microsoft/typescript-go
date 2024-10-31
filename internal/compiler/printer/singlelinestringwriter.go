package printer

import (
	"unicode/utf8"

	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/compiler/stringutil"
)

var SingleLineStringWriter EmitTextWriter = &singleLineStringWriter{}

type singleLineStringWriter struct {
	str string
}

// clear implements EmitTextWriter.
func (w *singleLineStringWriter) clear() {
	w.str = ""
}

// decreaseIndent implements EmitTextWriter.
func (w singleLineStringWriter) decreaseIndent() {
	// Do Nothing
}

// getColumn implements EmitTextWriter.
func (w singleLineStringWriter) getColumn() int {
	return 0
}

// getIndent implements EmitTextWriter.
func (w singleLineStringWriter) getIndent() int {
	return 0
}

// getLine implements EmitTextWriter.
func (w singleLineStringWriter) getLine() int {
	return 0
}

// getText implements EmitTextWriter.
func (w singleLineStringWriter) getText() string {
	return w.str
}

// getTextPos implements EmitTextWriter.
func (w singleLineStringWriter) getTextPos() int {
	return len(w.str)
}

// hasTrailingComment implements EmitTextWriter.
func (w singleLineStringWriter) hasTrailingComment() bool {
	return false
}

// hasTrailingWhitespace implements EmitTextWriter.
func (w singleLineStringWriter) hasTrailingWhitespace() bool {
	if len(w.str) == 0 {
		return false
	}
	ch, _ := utf8.DecodeLastRuneInString(w.str)
	if ch == utf8.RuneError {
		return false
	}
	return stringutil.IsWhiteSpaceLike(ch)
}

// increaseIndent implements EmitTextWriter.
func (w singleLineStringWriter) increaseIndent() {
	// Do Nothing
}

// isAtStartOfLine implements EmitTextWriter.
func (w singleLineStringWriter) isAtStartOfLine() bool {
	return false
}

// rawWrite implements EmitTextWriter.
func (w *singleLineStringWriter) rawWrite(s string) {
	w.str += s
}

// write implements EmitTextWriter.
func (w *singleLineStringWriter) write(s string) {
	w.str += s
}

// writeComment implements EmitTextWriter.
func (w *singleLineStringWriter) writeComment(text string) {
	w.str += text
}

// writeKeyword implements EmitTextWriter.
func (w *singleLineStringWriter) writeKeyword(text string) {
	w.str += text
}

// writeLine implements EmitTextWriter.
func (w *singleLineStringWriter) writeLine(force ...bool) {
	w.str += " "
}

// writeLiteral implements EmitTextWriter.
func (w *singleLineStringWriter) writeLiteral(s string) {
	w.str += s
}

// writeOperator implements EmitTextWriter.
func (w *singleLineStringWriter) writeOperator(text string) {
	w.str += text
}

// writeParameter implements EmitTextWriter.
func (w *singleLineStringWriter) writeParameter(text string) {
	w.str += text
}

// writeProperty implements EmitTextWriter.
func (w *singleLineStringWriter) writeProperty(text string) {
	w.str += text
}

// writePunctuation implements EmitTextWriter.
func (w *singleLineStringWriter) writePunctuation(text string) {
	w.str += text
}

// writeSpace implements EmitTextWriter.
func (w *singleLineStringWriter) writeSpace(text string) {
	w.str += text
}

// writeStringLiteral implements EmitTextWriter.
func (w *singleLineStringWriter) writeStringLiteral(text string) {
	w.str += text
}

// writeSymbol implements EmitTextWriter.
func (w *singleLineStringWriter) writeSymbol(text string, symbol compiler.Symbol) {
	w.str += text
}

// writeTrailingSemicolon implements EmitTextWriter.
func (w *singleLineStringWriter) writeTrailingSemicolon(text string) {
	w.str += text
}
