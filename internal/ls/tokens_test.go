package ls

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/parser"
	"github.com/microsoft/typescript-go/internal/repo"
	"github.com/microsoft/typescript-go/internal/scanner"
	"github.com/microsoft/typescript-go/internal/testutil/baseline"
	"github.com/microsoft/typescript-go/internal/testutil/jstest"
	"gotest.tools/v3/assert"
)

var testFiles = []string{
	filepath.Join(repo.TypeScriptSubmodulePath, "src/server/project.ts"),
}

func TestGetTokenAtPositionBaseline(t *testing.T) {
	t.Parallel()
	repo.SkipIfNoTypeScriptSubmodule(t)
	for _, fileName := range testFiles {
		fileText, err := os.ReadFile(fileName)
		assert.NilError(t, err)

		file := parser.ParseSourceFile("file.ts", string(fileText), core.ScriptTargetLatest, scanner.JSDocParsingModeParseAll)
		positions := make([]int, len(fileText))
		for i := range positions {
			positions[i] = i
		}

		// Get TypeScript tokens for all positions at once
		tsTokens := tsGetTokensAtPositions(t, string(fileText), positions)

		// Build the baseline output
		var output strings.Builder
		currentDiff := tokenDiff{}

		for pos, tsToken := range tsTokens {
			goToken := toTokenInfo(getTokenAtPosition(file, pos, true, false, nil))

			if goToken != tsToken {
				diff := tokenDiff{goToken: goToken, tsToken: tsToken}
				if currentDiff != diff {
					writeRangeDiff(&output, file, diff)
					currentDiff = diff
				}
			} else {
				currentDiff = tokenDiff{}
			}
		}
		if currentDiff != (tokenDiff{}) {
			writeRangeDiff(&output, file, currentDiff)
		}

		baseline.Run(t, filepath.Base(fileName)+".tokens.baseline.txt", output.String(), baseline.Options{
			Subfolder: "tokens",
		})
	}
}

type tokenDiff struct {
	goToken tokenInfo
	tsToken tokenInfo
}

type tokenInfo struct {
	Kind string `json:"kind"`
	Pos  int    `json:"pos"`
	End  int    `json:"end"`
}

func toTokenInfo(node *ast.Node) tokenInfo {
	kind := strings.Replace(node.Kind.String(), "Kind", "", 1)
	switch kind {
	case "EndOfFile":
		kind = "EndOfFileToken"
	}
	return tokenInfo{
		Kind: kind,
		Pos:  node.Pos(),
		End:  node.End(),
	}
}

func tsGetTokensAtPositions(t testing.TB, fileText string, positions []int) []tokenInfo {
	dir := t.TempDir()
	err := os.WriteFile(filepath.Join(dir, "file.ts"), []byte(fileText), 0o644)
	assert.NilError(t, err)
	script := `
		import fs from "fs";
		export default (ts, positions) => {
			positions = JSON.parse(positions);
			const fileText = fs.readFileSync("file.ts", "utf8");
			const file = ts.createSourceFile(
				"file.ts",
				fileText,
				{ languageVersion: ts.ScriptTarget.Latest, jsDocParsingMode: ts.JSDocParsingMode.ParseAll },
				/*setParentNodes*/ true
			);
			return positions.map(position => {
				const token = ts.getTokenAtPosition(file, position);
				return {
					kind: ts.Debug.formatSyntaxKind(token.kind),
					pos: token.pos,
					end: token.end,
				};
			});
		};`

	info, err := jstest.EvalNodeScriptWithTS[[]tokenInfo](t, script, dir, core.Must(core.StringifyJson(positions, "", "")))
	assert.NilError(t, err)
	return info
}

func tsGetTouchingPropertyName(t testing.TB, fileText string, positions []int) []tokenInfo {
	dir := t.TempDir()
	err := os.WriteFile(filepath.Join(dir, "file.ts"), []byte(fileText), 0o644)
	assert.NilError(t, err)
	script := `
		import fs from "fs";
		export default (ts, positions) => {
			positions = JSON.parse(positions);
			const fileText = fs.readFileSync("file.ts", "utf8");
			const file = ts.createSourceFile(
				"file.ts",
				fileText,
				{ languageVersion: ts.ScriptTarget.Latest, jsDocParsingMode: ts.JSDocParsingMode.ParseAll },
				/*setParentNodes*/ true
			);
			return positions.map(position => {
				const token = ts.getTouchingPropertyName(file, position);
				return {
					kind: ts.Debug.formatSyntaxKind(token.kind),
					pos: token.pos,
					end: token.end,
				};
			});
		};`

	info, err := jstest.EvalNodeScriptWithTS[[]tokenInfo](t, script, dir, core.Must(core.StringifyJson(positions, "", "")))
	assert.NilError(t, err)
	return info
}

func writeRangeDiff(output *strings.Builder, file *ast.SourceFile, diff tokenDiff) {
	lines := file.LineMap()
	tsStartLine, _ := core.PositionToLineAndCharacter(diff.tsToken.Pos, lines)
	tsEndLine, _ := core.PositionToLineAndCharacter(diff.tsToken.End, lines)
	goStartLine, _ := core.PositionToLineAndCharacter(diff.goToken.Pos, lines)
	goEndLine, _ := core.PositionToLineAndCharacter(diff.goToken.End, lines)
	contextLines := 2
	startLine := min(tsStartLine, goStartLine)
	endLine := max(tsEndLine, goEndLine)
	contextStart := max(0, startLine-contextLines)
	contextEnd := min(len(lines)-1, endLine+contextLines)

	if output.Len() > 0 {
		output.WriteString("\n\n")
	}

	output.WriteString(fmt.Sprintf("Line %d:\n", contextStart+1))
	output.WriteString(fmt.Sprintf("【TS: %s [%d, %d)】\n", diff.tsToken.Kind, diff.tsToken.Pos, diff.tsToken.End))
	output.WriteString(fmt.Sprintf("《Go: %s [%d, %d)》\n", diff.goToken.Kind, diff.goToken.Pos, diff.goToken.End))
	output.WriteString(file.Text[lines[contextStart]:lines[startLine]])
	for pos := int(lines[startLine]); pos < int(lines[endLine+1]); pos++ {
		if pos == diff.tsToken.Pos {
			output.WriteString("【")
		}
		if pos == diff.goToken.Pos {
			output.WriteString("《")
		}
		output.WriteByte(file.Text[pos])
		if pos == diff.tsToken.End-1 {
			output.WriteString("】")
		}
		if pos == diff.goToken.End-1 {
			output.WriteString("》")
		}
	}
	output.WriteString(file.Text[lines[endLine+1]:lines[contextEnd]])
}
