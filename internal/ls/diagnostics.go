package ls

import (
	"context"
	"slices"
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/diagnosticwriter"
	"github.com/microsoft/typescript-go/internal/ls/lsconv"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
)

func (l *LanguageService) ProvideDiagnostics(ctx context.Context, uri lsproto.DocumentUri) (lsproto.DocumentDiagnosticResponse, error) {
	program, file := l.getProgramAndFile(uri)

	syntactic := program.GetSyntacticDiagnostics(ctx, file)
	semantic := program.GetSemanticDiagnostics(ctx, file)
	suggestion := program.GetSuggestionDiagnostics(ctx, file)
	var declaration []*ast.Diagnostic
	if program.Options().GetEmitDeclarations() {
		declaration = program.GetDeclarationDiagnostics(ctx, file)
	}

	type diagsAndSeverity struct {
		diags    [][]*ast.Diagnostic
		severity lsproto.DiagnosticSeverity
	}

	// !!! user preference for suggestion diagnostics; keep only unnecessary/deprecated?
	// See: https://github.com/microsoft/vscode/blob/3dbc74129aaae102e5cb485b958fa5360e8d3e7a/extensions/typescript-language-features/src/languageFeatures/diagnostics.ts#L114
	// TODO: also implement reportStyleCheckAsWarnings to rewrite diags with Warning severity

	diags := []diagsAndSeverity{
		{diags: [][]*ast.Diagnostic{syntactic, semantic, declaration}, severity: lsproto.DiagnosticSeverityError},
		{diags: [][]*ast.Diagnostic{suggestion}, severity: lsproto.DiagnosticSeverityHint},
	}

	lspDiagnostics := make([]*lsproto.Diagnostic, 0, len(syntactic)+len(semantic)+len(suggestion)+len(declaration))

	for _, ds := range diags {
		for _, diagList := range ds.diags {
			for _, diag := range diagList {
				lspDiag := l.toLSPDiagnostic(ctx, diag, ds.severity)
				lspDiagnostics = append(lspDiagnostics, lspDiag)
			}
		}
	}

	return lsproto.RelatedFullDocumentDiagnosticReportOrUnchangedDocumentDiagnosticReport{
		FullDocumentDiagnosticReport: &lsproto.RelatedFullDocumentDiagnosticReport{
			Items: lspDiagnostics,
		},
	}, nil
}

func (l *LanguageService) toLSPDiagnostic(ctx context.Context, diagnostic *ast.Diagnostic, severity lsproto.DiagnosticSeverity) *lsproto.Diagnostic {
	clientOptions := lsproto.GetClientCapabilities(ctx).TextDocument.Diagnostic

	var relatedInformation []*lsproto.DiagnosticRelatedInformation
	if clientOptions.RelatedInformation {
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
	if len(clientOptions.TagSupport.ValueSet) > 0 && (diagnostic.ReportsUnnecessary() || diagnostic.ReportsDeprecated()) {
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
	diagnosticwriter.WriteFlattenedASTDiagnosticMessage(&b, diagnostic, "\n")
	return b.String()
}

func ptrToSliceIfNonEmpty[T any](s []T) *[]T {
	if len(s) == 0 {
		return nil
	}
	return &s
}
