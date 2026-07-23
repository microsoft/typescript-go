package lsconv

import (
	"context"
	"fmt"
	"net/url"
	"slices"
	"strings"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/debug"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/diagnosticwriter"
	"github.com/microsoft/typescript-go/internal/locale"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/spanmap"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type Converters struct {
	getLineMap       func(fileName string) *LSPLineMap
	positionEncoding lsproto.PositionEncodingKind
}

// Script is a source text the converters operate over. For a content-mapped file, Text() is the content
// mapper's transformed output and SpanMap() returns the map from that output back to the original text
// (OriginalText()); output ranges are then automatically converted to original coordinates (see
// ToLSPRange). For an ordinary file SpanMap() is nil and OriginalText() equals Text().
type Script interface {
	FileName() string
	Text() string
	SpanMap() *spanmap.SpanMap
	OriginalText() string
}

func NewConverters(positionEncoding lsproto.PositionEncodingKind, getLineMap func(fileName string) *LSPLineMap) *Converters {
	return &Converters{
		getLineMap:       getLineMap,
		positionEncoding: positionEncoding,
	}
}

// The To*/From* conversions map between an editor's coordinates and the coordinates the language service
// operates over. A content-mapped Script (SpanMap() != nil) is mapped through its span map; any other
// Script passes through unchanged.

func (c *Converters) ToLSPRange(script Script, textRange core.TextRange) (lsproto.Range, spanmap.Fidelity) {
	script, textRange, fidelity := mapOutputToOriginal(script, textRange)
	return lsproto.Range{
		Start: c.positionToLineAndCharacter(script, core.TextPos(textRange.Pos())),
		End:   c.positionToLineAndCharacter(script, core.TextPos(textRange.End())),
	}, fidelity
}

// ToLSPPosition is the single-position analog of ToLSPRange.
func (c *Converters) ToLSPPosition(script Script, position core.TextPos) (lsproto.Position, spanmap.Fidelity) {
	script, position, fidelity := mapOutputPositionToOriginal(script, position)
	return c.positionToLineAndCharacter(script, position), fidelity
}

func (c *Converters) ToLSPLocation(script Script, rng core.TextRange) (lsproto.Location, spanmap.Fidelity) {
	lspRange, fidelity := c.ToLSPRange(script, rng)
	return lsproto.Location{
		Uri:   FileNameToDocumentURI(script.FileName()),
		Range: lspRange,
	}, fidelity
}

// FromLSPRange converts an incoming LSP range to the offset range the language service operates over,
// mapping a content-mapped file's original text forward into its transformed text; it is the input analog
// of ToLSPRange.
func (c *Converters) FromLSPRange(script Script, textRange lsproto.Range, purpose spanmap.Purpose) []spanmap.MappedSpan {
	spans := script.SpanMap()
	if spans == nil {
		return []spanmap.MappedSpan{{
			Span: core.NewTextRange(
				int(c.lineAndCharacterToPosition(script, textRange.Start)),
				int(c.lineAndCharacterToPosition(script, textRange.End)),
			),
			Fidelity: spanmap.FidelityExact,
		}}
	}
	// A content-mapped script's line map is its original text's, so convert against that text and then map
	// the resulting original range forward into the transformed text.
	original := originalTextScript{fileName: script.FileName(), text: script.OriginalText()}
	origRange := core.NewTextRange(
		int(c.lineAndCharacterToPosition(original, textRange.Start)),
		int(c.lineAndCharacterToPosition(original, textRange.End)),
	)
	return spans.OriginalToGeneratedSpans(origRange, purpose)
}

// FromLSPPosition maps one incoming LSP position to every generated projection supporting purpose.
func (c *Converters) FromLSPPosition(script Script, position lsproto.Position, purpose spanmap.Purpose) []spanmap.MappedPosition {
	spans := script.SpanMap()
	if spans == nil {
		return []spanmap.MappedPosition{{Position: c.lineAndCharacterToPosition(script, position), Fidelity: spanmap.FidelityExact}}
	}
	original := originalTextScript{fileName: script.FileName(), text: script.OriginalText()}
	origOffset := c.lineAndCharacterToPosition(original, position)
	return spans.OriginalToGeneratedPositions(origOffset, purpose)
}

// mapOutputToOriginal maps a range in a content mapper's transformed output back to its original text,
// returning a Script over that original text and the mapping fidelity. Scripts that carry no span map are
// returned unchanged with FidelityExact.
func mapOutputToOriginal(script Script, textRange core.TextRange) (Script, core.TextRange, spanmap.Fidelity) {
	if script.SpanMap() == nil {
		return script, textRange, spanmap.FidelityExact
	}
	mapped, fidelity := script.SpanMap().GeneratedToOriginalSpan(textRange)
	return originalTextScript{fileName: script.FileName(), text: script.OriginalText()}, mapped, fidelity
}

// mapOutputPositionToOriginal is the single-position analog of mapOutputToOriginal.
func mapOutputPositionToOriginal(script Script, position core.TextPos) (Script, core.TextPos, spanmap.Fidelity) {
	if script.SpanMap() == nil {
		return script, position, spanmap.FidelityExact
	}
	mapped, fidelity := script.SpanMap().GeneratedToOriginalPosition(position)
	return originalTextScript{fileName: script.FileName(), text: script.OriginalText()}, mapped, fidelity
}

func LanguageKindToScriptKind(languageID lsproto.LanguageKind) core.ScriptKind {
	switch languageID {
	case "typescript":
		return core.ScriptKindTS
	case "typescriptreact":
		return core.ScriptKindTSX
	case "javascript":
		return core.ScriptKindJS
	case "javascriptreact":
		return core.ScriptKindJSX
	case "json":
		return core.ScriptKindJSON
	default:
		return core.ScriptKindUnknown
	}
}

// https://github.com/microsoft/vscode-uri/blob/edfdccd976efaf4bb8fdeca87e97c47257721729/src/uri.ts#L455
var extraEscapeReplacer = strings.NewReplacer(
	":", "%3A",
	"/", "%2F",
	"?", "%3F",
	"#", "%23",
	"[", "%5B",
	"]", "%5D",
	"@", "%40",

	"!", "%21",
	"$", "%24",
	"&", "%26",
	"'", "%27",
	"(", "%28",
	")", "%29",
	"*", "%2A",
	"+", "%2B",
	",", "%2C",
	";", "%3B",
	"=", "%3D",

	" ", "%20",
)

func FileNameToDocumentURI(fileName string) lsproto.DocumentUri {
	if bundled.IsBundled(fileName) {
		return lsproto.DocumentUri(fileName)
	}
	if tspath.IsDynamicFileName(fileName) {
		scheme, rest, ok := strings.Cut(fileName[2:], "/")
		if !ok {
			panic("invalid file name: " + fileName)
		}
		authority, path, ok := strings.Cut(rest, "/")
		if !ok {
			panic("invalid file name: " + fileName)
		}
		if authority == "ts-nul-authority" {
			return lsproto.DocumentUri(scheme + ":" + path)
		}
		return lsproto.DocumentUri(scheme + "://" + authority + "/" + path)
	}

	volume, fileName, _ := tspath.SplitVolumePath(fileName)
	if volume != "" {
		volume = "/" + extraEscapeReplacer.Replace(volume)
	}

	fileName = strings.TrimPrefix(fileName, "//")

	parts := strings.Split(fileName, "/")
	for i, part := range parts {
		parts[i] = extraEscapeReplacer.Replace(url.PathEscape(part))
	}

	return lsproto.DocumentUri("file://" + volume + strings.Join(parts, "/"))
}

func (c *Converters) lineAndCharacterToPosition(script Script, lineAndCharacter lsproto.Position) core.TextPos {
	// UTF-8/16 0-indexed line and character to UTF-8 offset
	debug.Assert(script.SpanMap() == nil, "raw coordinate conversion requires a non-content-mapped script")

	lineMap := c.getLineMap(script.FileName())

	line := core.TextPos(lineAndCharacter.Line)
	char := core.TextPos(lineAndCharacter.Character)

	textLen := core.TextPos(len(script.Text()))

	// Clamp line to valid range.
	if int(line) >= len(lineMap.LineStarts) {
		return textLen
	}

	start := lineMap.LineStarts[line]

	// Determine the end of this line (start of next line, or end of text).
	var lineEnd core.TextPos
	if int(line)+1 < len(lineMap.LineStarts) {
		lineEnd = lineMap.LineStarts[int(line)+1]
	} else {
		lineEnd = textLen
	}

	if lineMap.AsciiOnly || c.positionEncoding == lsproto.PositionEncodingKindUTF8 {
		return max(start, min(start+char, lineEnd))
	}

	// Scan from line start counting UTF-16 code units to find the byte position.
	// Uses DecodeRuneInString (not range + RuneLen) so that invalid UTF-8 bytes
	// advance by their actual size (1) rather than RuneLen(RuneError) == 3.
	// This matches the approach in scanner.ComputePositionOfLineAndUTF16Character.
	var utf16Char core.TextPos
	pos := int(start)
	end := int(lineEnd)
	text := script.Text()
	for pos < end {
		r, size := utf8.DecodeRuneInString(text[pos:])
		u16Len := core.TextPos(utf16.RuneLen(r))
		if utf16Char+u16Len > char {
			break
		}
		utf16Char += u16Len
		pos += size
	}

	return core.TextPos(pos)
}

func (c *Converters) positionToLineAndCharacter(script Script, position core.TextPos) lsproto.Position {
	// UTF-8 offset to UTF-8/16 0-indexed line and character
	debug.Assert(script.SpanMap() == nil, "raw coordinate conversion requires a non-content-mapped script")

	position = max(0, min(position, core.TextPos(len(script.Text()))))

	lineMap := c.getLineMap(script.FileName())

	line, isLineStart := slices.BinarySearch(lineMap.LineStarts, position)
	if !isLineStart {
		line--
	}
	line = max(0, min(line, len(lineMap.LineStarts)-1))

	// The current line ranges from lineMap.LineStarts[line] (or 0) to lineMap.LineStarts[line+1] (or len(text)).

	start := lineMap.LineStarts[line]

	var character core.TextPos
	if lineMap.AsciiOnly || c.positionEncoding == lsproto.PositionEncodingKindUTF8 {
		character = position - start
	} else {
		// We need to rescan the text as UTF-16 to find the character offset.
		for _, r := range script.Text()[start:position] {
			character += core.TextPos(utf16.RuneLen(r))
		}
	}

	return lsproto.Position{
		Line:      uint32(line),
		Character: uint32(character),
	}
}

type diagnosticOptions struct {
	reportStyleChecksAsWarnings bool
	relatedInformation          bool
	tagValueSet                 []lsproto.DiagnosticTag
	visualStudio                bool
}

// DiagnosticToLSPPull converts a diagnostic for pull diagnostics (textDocument/diagnostic)
func DiagnosticToLSPPull(ctx context.Context, converters *Converters, diagnostic *ast.Diagnostic, reportStyleChecksAsWarnings bool) *lsproto.Diagnostic {
	clientCaps := lsproto.GetClientCapabilities(ctx)
	clientDiagnosticCaps := clientCaps.TextDocument.Diagnostic
	return diagnosticToLSP(ctx, converters, diagnostic, diagnosticOptions{
		reportStyleChecksAsWarnings: reportStyleChecksAsWarnings, // !!! get through context UserPreferences
		relatedInformation:          clientDiagnosticCaps.RelatedInformation,
		tagValueSet:                 clientDiagnosticCaps.TagSupport.ValueSet,
		visualStudio:                clientCaps.VSSupportsVisualStudioExtensions,
	})
}

// DiagnosticToLSPPush converts a diagnostic for push diagnostics (textDocument/publishDiagnostics)
func DiagnosticToLSPPush(ctx context.Context, converters *Converters, diagnostic *ast.Diagnostic) *lsproto.Diagnostic {
	clientCaps := lsproto.GetClientCapabilities(ctx)
	clientDiagnosticCaps := clientCaps.TextDocument.PublishDiagnostics
	return diagnosticToLSP(ctx, converters, diagnostic, diagnosticOptions{
		relatedInformation: clientDiagnosticCaps.RelatedInformation,
		tagValueSet:        clientDiagnosticCaps.TagSupport.ValueSet,
		visualStudio:       clientCaps.VSSupportsVisualStudioExtensions,
	})
}

// https://github.com/microsoft/vscode/blob/93e08afe0469712706ca4e268f778cfadf1a43ef/extensions/typescript-language-features/src/typeScriptServiceClientHost.ts#L40C7-L40C29
var styleCheckDiagnostics = collections.NewSetFromItems(
	diagnostics.X_0_is_declared_but_never_used.Code(),
	diagnostics.X_0_is_declared_but_its_value_is_never_read.Code(),
	diagnostics.Property_0_is_declared_but_its_value_is_never_read.Code(),
	diagnostics.All_imports_in_import_declaration_are_unused.Code(),
	diagnostics.Unreachable_code_detected.Code(),
	diagnostics.Unused_label.Code(),
	diagnostics.Fallthrough_case_in_switch.Code(),
	diagnostics.Not_all_code_paths_return_a_value.Code(),
)

func diagnosticToLSP(ctx context.Context, converters *Converters, diagnostic *ast.Diagnostic, opts diagnosticOptions) *lsproto.Diagnostic {
	locale := locale.FromContext(ctx)
	severity := diagnosticSeverity(diagnostic.Category())

	if opts.reportStyleChecksAsWarnings && severity == lsproto.DiagnosticSeverityError && styleCheckDiagnostics.Has(diagnostic.Code()) {
		severity = lsproto.DiagnosticSeverityWarning
	}

	var relatedInformation []*lsproto.DiagnosticRelatedInformation
	if opts.relatedInformation {
		relatedInformation = make([]*lsproto.DiagnosticRelatedInformation, 0, len(diagnostic.RelatedInformation()))
		for _, related := range diagnostic.RelatedInformation() {
			script, loc := diagnosticScriptAndRange(related.File(), related.Loc(), related.Source())
			relatedRange, fidelity := converters.ToLSPRange(script, loc)
			if fidelity.IsNone() {
				// Related diagnostic information cannot omit its location. Use an explicit file-level
				// location instead of presenting the synthesized span's insertion point as related source.
				relatedRange = lsproto.Range{}
			}
			relatedInformation = append(relatedInformation, &lsproto.DiagnosticRelatedInformation{
				Location: lsproto.Location{
					Uri:   FileNameToDocumentURI(related.File().FileName()),
					Range: relatedRange,
				},
				Message: related.Localize(locale),
			})
		}
	}

	var tags []lsproto.DiagnosticTag
	if len(opts.tagValueSet) > 0 && (diagnostic.ReportsUnnecessary() || diagnostic.ReportsDeprecated()) {
		tags = make([]lsproto.DiagnosticTag, 0, 2)
		if diagnostic.ReportsUnnecessary() && slices.Contains(opts.tagValueSet, lsproto.DiagnosticTagUnnecessary) {
			tags = append(tags, lsproto.DiagnosticTagUnnecessary)
		}
		if diagnostic.ReportsDeprecated() && slices.Contains(opts.tagValueSet, lsproto.DiagnosticTagDeprecated) {
			tags = append(tags, lsproto.DiagnosticTagDeprecated)
		}
	}

	// For diagnostics without a file (e.g., program diagnostics), use a zero range
	var lspRange lsproto.Range
	if diagnostic.File() != nil {
		script, loc := diagnosticScriptAndRange(diagnostic.File(), diagnostic.Loc(), diagnostic.Source())
		var fidelity spanmap.Fidelity
		lspRange, fidelity = converters.ToLSPRange(script, loc)
		if fidelity.IsNone() {
			// Diagnostics must carry a range. A zero range honestly means "this file" when the
			// diagnostic arose entirely in synthesized code and has no original source span.
			lspRange = lsproto.Range{}
		}
	}

	var code *lsproto.IntegerOrString
	sourceText := diagnostic.Source()
	if sourceText == "" {
		sourceText = "ts"
	}
	if opts.visualStudio {
		code = &lsproto.IntegerOrString{
			String: new(fmt.Sprintf("TS%d", diagnostic.Code())),
		}
	} else {
		code = &lsproto.IntegerOrString{
			Integer: new(diagnostic.Code()),
		}
	}

	return &lsproto.Diagnostic{
		Range:              lspRange,
		Code:               code,
		Severity:           &severity,
		Message:            lsproto.StringOrMarkupContent{String: new(messageChainToString(diagnostic, locale))},
		Source:             &sourceText,
		RelatedInformation: ptrToSliceIfNonEmpty(relatedInformation),
		Tags:               ptrToSliceIfNonEmpty(tags),
	}
}

// diagnosticScriptAndRange resolves the text basis and range to report a diagnostic against. For a
// content-mapped file it maps the diagnostic's transformed range back to the original text so
// the range lines up with what the editor shows; the original text's line map is already what
// getLineMap returns for the file. A range in synthesized code has no original counterpart, so it is
// surfaced at the top of the file. Non-mapped files are returned unchanged.
func diagnosticScriptAndRange(file *ast.SourceFile, loc core.TextRange, source string) (Script, core.TextRange) {
	if file == nil || file.SpanMap() == nil {
		return file, loc
	}
	original := originalTextScript{fileName: file.FileName(), text: file.OriginalText()}
	if source != "" {
		// A content mapper's own diagnostics already carry original-text ranges.
		return original, loc
	}
	mapped, fidelity := file.SpanMap().GeneratedToOriginalSpan(loc)
	if fidelity == spanmap.FidelityNone {
		// Entirely synthesized code has no original location; surface it at the top of the file.
		return original, core.NewTextRange(0, 0)
	}
	return original, mapped
}

// originalTextScript presents a content-mapped file's original (untransformed) text as a Script, so that
// ranges already mapped into that text convert to the correct line/character positions.
type originalTextScript struct {
	fileName string
	text     string
}

func (s originalTextScript) FileName() string        { return s.fileName }
func (s originalTextScript) Text() string            { return s.text }
func (s originalTextScript) OriginalText() string    { return s.text }
func (originalTextScript) SpanMap() *spanmap.SpanMap { return nil }

// diagnosticSeverity maps a diagnostic category to its LSP severity.
func diagnosticSeverity(category diagnostics.Category) lsproto.DiagnosticSeverity {
	switch category {
	case diagnostics.CategorySuggestion:
		return lsproto.DiagnosticSeverityHint
	case diagnostics.CategoryMessage:
		return lsproto.DiagnosticSeverityInformation
	case diagnostics.CategoryWarning:
		return lsproto.DiagnosticSeverityWarning
	default:
		return lsproto.DiagnosticSeverityError
	}
}

func messageChainToString(diagnostic *ast.Diagnostic, locale locale.Locale) string {
	if len(diagnostic.MessageChain()) == 0 {
		return diagnostic.Localize(locale)
	}
	var b strings.Builder
	diagnosticwriter.WriteFlattenedASTDiagnosticMessage(&b, diagnostic, "\n", locale)
	return b.String()
}

func ptrToSliceIfNonEmpty[T any](s []T) *[]T {
	if len(s) == 0 {
		return nil
	}
	return &s
}
