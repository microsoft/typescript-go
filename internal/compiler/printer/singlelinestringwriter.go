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

func (w *singleLineStringWriter) clear() {
	w.str = ""
}

func (w singleLineStringWriter) decreaseIndent() {
	// Do Nothing
}

func (w singleLineStringWriter) getColumn() int {
	return 0
}

func (w singleLineStringWriter) getIndent() int {
	return 0
}

func (w singleLineStringWriter) getLine() int {
	return 0
}

func (w singleLineStringWriter) getText() string {
	return w.str
}

func (w singleLineStringWriter) getTextPos() int {
	return len(w.str)
}

func (w singleLineStringWriter) hasTrailingComment() bool {
	return false
}

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

func (w singleLineStringWriter) increaseIndent() {
	// Do Nothing
}

func (w singleLineStringWriter) isAtStartOfLine() bool {
	return false
}

func (w *singleLineStringWriter) rawWrite(s string) {
	w.str += s
}

func (w *singleLineStringWriter) write(s string) {
	w.str += s
}

func (w *singleLineStringWriter) writeComment(text string) {
	w.str += text
}

func (w *singleLineStringWriter) writeKeyword(text string) {
	w.str += text
}

func (w *singleLineStringWriter) writeLine(force ...bool) {
	w.str += " "
}

func (w *singleLineStringWriter) writeLiteral(s string) {
	w.str += s
}

func (w *singleLineStringWriter) writeOperator(text string) {
	w.str += text
}

func (w *singleLineStringWriter) writeParameter(text string) {
	w.str += text
}

func (w *singleLineStringWriter) writeProperty(text string) {
	w.str += text
}

func (w *singleLineStringWriter) writePunctuation(text string) {
	w.str += text
}

func (w *singleLineStringWriter) writeSpace(text string) {
	w.str += text
}

func (w *singleLineStringWriter) writeStringLiteral(text string) {
	w.str += text
}

func (w *singleLineStringWriter) writeSymbol(text string, symbol compiler.Symbol) {
	w.str += text
}

func (w *singleLineStringWriter) writeTrailingSemicolon(text string) {
	w.str += text
}
