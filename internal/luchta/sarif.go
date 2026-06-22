package luchta

import (
	"encoding/json"
	"strconv"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/diagnosticwriter"
	"github.com/microsoft/typescript-go/internal/locale"
	"github.com/microsoft/typescript-go/internal/scanner"
)

// sarifMimeType is the MIME type luchta dispatches on to pretty-print a worker
// report as SARIF (luchta crates/luchta-cli/src/reports). Dispatch is by MIME
// only — never filename or extension.
const sarifMimeType = "application/sarif+json"

// sarifReportFilename is the report basename attached to the task result. It is
// used only by `luchta logs --file`; native rendering keys off the MIME type.
const sarifReportFilename = "tsc.sarif"

// Minimal SARIF 2.1.0 document, covering the fields luchta consumes plus enough
// structure to be a valid log for other SARIF tooling.
type sarifLog struct {
	Schema  string     `json:"$schema"`
	Version string     `json:"version"`
	Runs    []sarifRun `json:"runs"`
}

type sarifRun struct {
	Tool    sarifTool     `json:"tool"`
	Results []sarifResult `json:"results"`
}

type sarifTool struct {
	Driver sarifDriver `json:"driver"`
}

type sarifDriver struct {
	Name           string `json:"name"`
	InformationURI string `json:"informationUri,omitempty"`
}

type sarifResult struct {
	RuleID    string          `json:"ruleId,omitempty"`
	Level     string          `json:"level"`
	Message   sarifMessage    `json:"message"`
	Locations []sarifLocation `json:"locations,omitempty"`
}

type sarifMessage struct {
	Text string `json:"text"`
}

type sarifLocation struct {
	PhysicalLocation sarifPhysicalLocation `json:"physicalLocation"`
}

type sarifPhysicalLocation struct {
	ArtifactLocation sarifArtifactLocation `json:"artifactLocation"`
	Region           *sarifRegion          `json:"region,omitempty"`
}

type sarifArtifactLocation struct {
	URI string `json:"uri"`
}

type sarifRegion struct {
	StartLine   int `json:"startLine"`
	StartColumn int `json:"startColumn"`
	EndLine     int `json:"endLine,omitempty"`
	EndColumn   int `json:"endColumn,omitempty"`
}

// DiagnosticsToSARIF renders compiler diagnostics as a SARIF 2.1.0 log string,
// suitable for emitting as a worker report with mime type [sarifMimeType].
func DiagnosticsToSARIF(diags []*ast.Diagnostic) string {
	results := make([]sarifResult, 0, len(diags))
	for _, d := range diags {
		results = append(results, diagnosticToSARIFResult(d))
	}
	doc := sarifLog{
		Schema:  "https://json.schemastore.org/sarif-2.1.0.json",
		Version: "2.1.0",
		Runs: []sarifRun{{
			Tool: sarifTool{Driver: sarifDriver{
				Name:           "luchta-tsc-worker",
				InformationURI: "https://github.com/microsoft/typescript-go",
			}},
			Results: results,
		}},
	}
	b, err := json.Marshal(doc)
	if err != nil {
		// Diagnostics are plain data, so marshaling cannot realistically fail;
		// fall back to a valid empty log rather than emitting malformed JSON.
		return `{"$schema":"https://json.schemastore.org/sarif-2.1.0.json","version":"2.1.0","runs":[]}`
	}
	return string(b)
}

func diagnosticToSARIFResult(d *ast.Diagnostic) sarifResult {
	r := sarifResult{
		RuleID:  diagnosticRuleID(d),
		Level:   sarifLevel(d.Category()),
		Message: sarifMessage{Text: diagnosticwriter.FlattenDiagnosticMessage(diagnosticwriter.WrapASTDiagnostic(d), "\n", locale.Default)},
	}
	if file := d.File(); file != nil {
		// SARIF lines/columns are 1-based; tsgo's are 0-based.
		startLine, startCol := scanner.GetECMALineAndUTF16CharacterOfPosition(file, d.Pos())
		endLine, endCol := scanner.GetECMALineAndUTF16CharacterOfPosition(file, d.End())
		r.Locations = []sarifLocation{{
			PhysicalLocation: sarifPhysicalLocation{
				// Absolute path so the link is clickable regardless of the
				// directory `luchta logs` is run from. tsgo normalizes file
				// names to absolute paths.
				ArtifactLocation: sarifArtifactLocation{URI: file.FileName()},
				Region: &sarifRegion{
					StartLine:   startLine + 1,
					StartColumn: int(startCol) + 1,
					EndLine:     endLine + 1,
					EndColumn:   int(endCol) + 1,
				},
			},
		}}
	}
	return r
}

func diagnosticRuleID(d *ast.Diagnostic) string {
	if d.Code() == 0 {
		return ""
	}
	return "TS" + strconv.Itoa(int(d.Code()))
}

func sarifLevel(c diagnostics.Category) string {
	switch c {
	case diagnostics.CategoryError:
		return "error"
	case diagnostics.CategoryWarning:
		return "warning"
	default: // CategorySuggestion, CategoryMessage
		return "note"
	}
}
