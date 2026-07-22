package lsproto

import (
	"cmp"
	"fmt"
	"strings"
)

// Implements a cmp.Compare like function for two Position
// ComparePositions(pos, other) == cmp.Compare(pos, other)
func ComparePositions(pos, other Position) int {
	if lineComp := cmp.Compare(pos.Line, other.Line); lineComp != 0 {
		return lineComp
	}
	return cmp.Compare(pos.Character, other.Character)
}

// Implements a cmp.Compare like function for two Range
// CompareRanges(lsRange, other) == cmp.Compare(lsRange, other)
//
//	Range.Start is compared before Range.End
func CompareRanges(lsRange, other Range) int {
	if startComp := ComparePositions(lsRange.Start, other.Start); startComp != 0 {
		return startComp
	}
	return ComparePositions(lsRange.End, other.End)
}

// AsString returns the plain text of a StringOrMarkupContent, reading the
// MarkupContent value when the message is not a plain string.
func (m StringOrMarkupContent) AsString() string {
	if m.String != nil {
		return *m.String
	}
	if m.MarkupContent != nil {
		return m.MarkupContent.Value
	}
	return ""
}

func (m IntegerOrString) AsString() string {
	codeStr := ""
	if m.String != nil {
		codeStr = *m.String
	} else if m.Integer != nil {
		codeStr = fmt.Sprintf("%v", *m.Integer)
	} else {
		codeStr = "-1"
	}
	return codeStr
}

func diagnosticExistsInSlice(elem *Diagnostic, diags []*Diagnostic) bool {
	for _, diag := range diags {
		if diagnosticsEqual(elem, diag) {
			return true
		}
	}
	return false
}

func diagnosticsEqual(diag1 *Diagnostic, diag2 *Diagnostic) bool {
	if diagnosticCodesEqual(diag1.Code, diag2.Code) && diagnosticMessagesEqual(diag1.Message, diag2.Message) && CompareRanges(diag1.Range, diag2.Range) == 0 {
		return true
	}
	return false
}

func diagnosticCodesEqual(code1 *IntegerOrString, code2 *IntegerOrString) bool {
	if code1.String != nil && code2.String != nil {
		return *code1.String == *code2.String
	}
	if code1.Integer != nil && code2.Integer != nil {
		return *code1.Integer == *code2.Integer
	}
	return false
}

func diagnosticMessagesEqual(message1 StringOrMarkupContent, message2 StringOrMarkupContent) bool {
	if message1.String != nil && message2.String != nil {
		return *message1.String == *message2.String
	}
	if message1.MarkupContent != nil && message2.MarkupContent != nil {
		return message1.MarkupContent.Kind == message2.MarkupContent.Kind && message1.MarkupContent.Value == message2.MarkupContent.Value
	}
	return false
}

func CompareDiagnostics(preEmit []*Diagnostic, postEmit []*Diagnostic, sanitizedDiffEnabled bool) (string, string) {
	sanitizedMsgBuilder := strings.Builder{}
	msgBuilder := strings.Builder{}
	for _, elem := range preEmit {
		if !diagnosticExistsInSlice(elem, postEmit) {
			msgBuilder.WriteString(fmt.Sprintf("Diagnostic `%v` was present before emit but not after emit\n", elem.AsString()))
			if sanitizedDiffEnabled {
				sanitizedMsgBuilder.WriteString(fmt.Sprintf("Diagnostic %v was present before emit but not after emit\n", elem.CodeAsString()))
			}
		}
	}
	for _, elem := range postEmit {
		if !diagnosticExistsInSlice(elem, preEmit) {
			msgBuilder.WriteString(fmt.Sprintf("Diagnostic `%v` was present after emit but not before emit\n", elem.AsString()))
			if sanitizedDiffEnabled {
				sanitizedMsgBuilder.WriteString(fmt.Sprintf("Diagnostic %v was present after emit but not before emit\n", elem.CodeAsString()))
			}
		}
	}
	return msgBuilder.String(), sanitizedMsgBuilder.String()
}

func (elem *Diagnostic) AsString() string {
	return fmt.Sprintf("%v (%v:%v-%v:%v): %v", elem.Code.AsString(), elem.Range.Start.Line, elem.Range.Start.Character, elem.Range.End.Line, elem.Range.End.Character, elem.Message.AsString())
}

func (elem *Diagnostic) CodeAsString() string {
	return fmt.Sprintf("Code(%v)", elem.Code.AsString())
}
