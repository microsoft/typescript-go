package ls

import (
	"context"
	"slices"
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/diagnosticwriter"
	"github.com/microsoft/typescript-go/internal/ls/lsconv"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
)

func (l *LanguageService) ProvideDiagnostics(ctx context.Context, uri lsproto.DocumentUri, clientOptions *lsproto.DiagnosticClientCapabilities) (lsproto.DocumentDiagnosticResponse, error) {
	program, file := l.getProgramAndFile(uri)

	diagnostics := make([][]*ast.Diagnostic, 0, 4)
	diagnostics = append(diagnostics, program.GetSyntacticDiagnostics(ctx, file))
	diagnostics = append(diagnostics, program.GetSemanticDiagnostics(ctx, file))
	// !!! user preference for suggestion diagnostics; keep only unnecessary/deprecated?
	// See: https://github.com/microsoft/vscode/blob/3dbc74129aaae102e5cb485b958fa5360e8d3e7a/extensions/typescript-language-features/src/languageFeatures/diagnostics.ts#L114
	diagnostics = append(diagnostics, program.GetSuggestionDiagnostics(ctx, file))
	if program.Options().GetEmitDeclarations() {
		diagnostics = append(diagnostics, program.GetDeclarationDiagnostics(ctx, file))
	}

	return lsproto.RelatedFullDocumentDiagnosticReportOrUnchangedDocumentDiagnosticReport{
		FullDocumentDiagnosticReport: &lsproto.RelatedFullDocumentDiagnosticReport{
			Items: l.toLSPDiagnostics(clientOptions, diagnostics...),
		},
	}, nil
}

func (l *LanguageService) toLSPDiagnostics(clientOptions *lsproto.DiagnosticClientCapabilities, diagnostics ...[]*ast.Diagnostic) []*lsproto.Diagnostic {
	size := 0
	for _, diagSlice := range diagnostics {
		size += len(diagSlice)
	}
	lspDiagnostics := make([]*lsproto.Diagnostic, 0, size)
	for _, diagSlice := range diagnostics {
		for _, diag := range diagSlice {
			lspDiagnostics = append(lspDiagnostics, l.toLSPDiagnostic(clientOptions, diag))
		}
	}
	return lspDiagnostics
}

func (l *LanguageService) toLSPDiagnostic(clientOptions *lsproto.DiagnosticClientCapabilities, diagnostic *ast.Diagnostic) *lsproto.Diagnostic {
	var severity lsproto.DiagnosticSeverity
	switch diagnostic.Category() {
	case diagnostics.CategorySuggestion:
		severity = lsproto.DiagnosticSeverityHint
	case diagnostics.CategoryMessage:
		severity = lsproto.DiagnosticSeverityInformation
	case diagnostics.CategoryWarning:
		severity = lsproto.DiagnosticSeverityWarning
	default:
		severity = lsproto.DiagnosticSeverityError
	}

	var relatedInformation []*lsproto.DiagnosticRelatedInformation
	if clientOptions != nil && ptrIsTrue(clientOptions.RelatedInformation) {
		relatedInformation = make([]*lsproto.DiagnosticRelatedInformation, 0, len(diagnostic.RelatedInformation()))
		for _, related := range diagnostic.RelatedInformation() {
			relatedInformation = append(relatedInformation, &lsproto.DiagnosticRelatedInformation{
				Location: lsproto.Location{
					Uri:   lsconv.FileNameToDocumentURI(related.File().FileName()),
					Range: l.converters.ToLSPRange(related.File(), related.Loc()),
				},
				Message: related.Message(),
			})
		}
	}

	var tags []lsproto.DiagnosticTag
	if clientOptions != nil && clientOptions.TagSupport != nil && (diagnostic.ReportsUnnecessary() || diagnostic.ReportsDeprecated()) {
		tags = make([]lsproto.DiagnosticTag, 0, 2)
		if diagnostic.ReportsUnnecessary() && slices.Contains(clientOptions.TagSupport.ValueSet, lsproto.DiagnosticTagUnnecessary) {
			tags = append(tags, lsproto.DiagnosticTagUnnecessary)
		}
		if diagnostic.ReportsDeprecated() && slices.Contains(clientOptions.TagSupport.ValueSet, lsproto.DiagnosticTagDeprecated) {
			tags = append(tags, lsproto.DiagnosticTagDeprecated)
		}
	}

	return &lsproto.Diagnostic{
		Range: l.converters.ToLSPRange(diagnostic.File(), diagnostic.Loc()),
		Code: &lsproto.IntegerOrString{
			Integer: ptrTo(diagnostic.Code()),
		},
		Severity:           &severity,
		Message:            messageChainToString(diagnostic),
		Source:             ptrTo("ts"),
		RelatedInformation: ptrToSliceIfNonEmpty(relatedInformation),
		Tags:               ptrToSliceIfNonEmpty(tags),
	}
}

func messageChainToString(diagnostic *ast.Diagnostic) string {
	if len(diagnostic.MessageChain()) == 0 {
		return diagnostic.Message()
	}
	var b strings.Builder
	diagnosticwriter.WriteFlattenedDiagnosticMessage(&b, diagnostic, "\n")
	return b.String()
}

func ptrToSliceIfNonEmpty[T any](s []T) *[]T {
	if len(s) == 0 {
		return nil
	}
	return &s
}
