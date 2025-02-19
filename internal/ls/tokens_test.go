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
	"github.com/microsoft/typescript-go/internal/testutil/jstest"
	"gotest.tools/v3/assert"
)

var testFiles = []string{
	filepath.Join(repo.TypeScriptSubmodulePath, "src/server/project.ts"),
}

func BenchmarkTokens(b *testing.B) {
	repo.SkipIfNoTypeScriptSubmodule(b)
	for _, fileName := range testFiles {
		fileText, err := os.ReadFile(fileName)
		assert.NilError(b, err)
		positionCount := 50
		positions := make([]int, positionCount)
		for i := range positionCount {
			positions[i] = i * len(fileText) / positionCount
		}
		file := parser.ParseSourceFile("file.ts", string(fileText), core.ScriptTargetLatest, scanner.JSDocParsingModeParseAll)
		ast.SetParentInChildren(file.AsNode())
		for _, pos := range positions {
			b.Run(fmt.Sprintf("getTokenAtPosition:%s:%d", filepath.Base(fileName), pos), func(b *testing.B) {
				for range b.N {
					getTokenAtPosition(file, pos, true /*allowPositionInLeadingTrivia*/, false /*includeEndPosition*/, nil)
				}
			})
			b.Run(fmt.Sprintf("getTokenAtPosition_fast:%s:%d", filepath.Base(fileName), pos), func(b *testing.B) {
				for range b.N {
					getTokenAtPosition_fast(file, pos, true /*allowPositionInLeadingTrivia*/, false /*includeEndPosition*/, nil)
				}
			})
		}
	}
}

func TestGetTokenAtPositionFast(t *testing.T) {
	t.Parallel()
	repo.SkipIfNoTypeScriptSubmodule(t)
	for _, fileName := range testFiles {
		fileText, err := os.ReadFile(fileName)
		assert.NilError(t, err)
		positionCount := len(fileText)
		positions := make([]int, positionCount)
		for i := range positionCount {
			positions[i] = i * len(fileText) / positionCount
		}
		file := parser.ParseSourceFile("file.ts", string(fileText), core.ScriptTargetLatest, scanner.JSDocParsingModeParseAll)
		for _, pos := range positions {
			t.Run(fmt.Sprintf("pos: %d", pos), func(t *testing.T) {
				t.Parallel()
				slow := getTokenAtPosition(file, pos, true /*allowPositionInLeadingTrivia*/, false /*includeEndPosition*/, nil)
				fast := getTokenAtPosition_fast(file, pos, true, false, nil)
				if fast.Kind == ast.KindJSDoc && slow.Kind != ast.KindJSDoc && pos < fast.Pos() {
					// JSDoc positions are incorrect, which has been worked around in getTokenAtPosition_fast.
					// When the cursor is in whitespace before JSDoc, the fast version will return the JSDoc token,
					// whereas the slow version will return the node containing the JSDoc, whose first non-comment
					// token is after the JSDoc.
					return
				}
				assert.Equal(
					t,
					slow.Kind,
					fast.Kind,
					pos,
				)
			})
		}
	}
}

func TestGetTokenAtPosition(t *testing.T) {
	t.Parallel()
	jstest.SkipIfNoNodeJS(t)
	repo.SkipIfNoTypeScriptSubmodule(t)
	for _, fileName := range testFiles {
		fileText, err := os.ReadFile(fileName)
		assert.NilError(t, err)
		file := parser.ParseSourceFile("file.ts", string(fileText), core.ScriptTargetLatest, scanner.JSDocParsingModeParseAll)
		ast.SetParentInChildren(file.AsNode())
		positionCount := 100
		positions := make([]int, positionCount)
		for i := range positionCount {
			positions[i] = i * len(fileText) / positionCount
		}

		tsTokens := tsGetTokensAtPositions(t, string(fileText), positions)
		t.Run(fileName, func(t *testing.T) {
			t.Parallel()
			for i, pos := range positions {
				t.Run(fmt.Sprintf("pos: %d", pos), func(t *testing.T) {
					t.Parallel()
					goToken := getTokenAtPosition(file, pos, true /*allowPositionInLeadingTrivia*/, false /*includeEndPosition*/, nil)
					if goToken.Kind == ast.KindJSDocText && strings.HasPrefix(tsTokens[i].Kind, "JSDoc") {
						// Strada sometimes stored plain-text JSDoc comments as strings
						// on JSDoc nodes, whereas Corsa stores them as JSDocText nodes.
						// It's fine for Corsa to return a deeper, more specific node in
						// this case.
						return
					}
					assert.Equal(t, tsTokens[i], toTokenInfo(goToken))
				})
			}
		})
	}
}

func TestGetTouchingPropertyNameFast(t *testing.T) {
	t.Parallel()
	repo.SkipIfNoTypeScriptSubmodule(t)
	for _, fileName := range testFiles {
		fileText, err := os.ReadFile(fileName)
		assert.NilError(t, err)
		positionCount := len(fileText)
		positions := make([]int, positionCount)
		for i := range positionCount {
			positions[i] = i * len(fileText) / positionCount
		}
		file := parser.ParseSourceFile("file.ts", string(fileText), core.ScriptTargetLatest, scanner.JSDocParsingModeParseAll)
		for _, pos := range positions {
			t.Run(fmt.Sprintf("pos: %d", pos), func(t *testing.T) {
				t.Parallel()
				slow := getTouchingPropertyName(file, pos)
				if slow.Kind == ast.KindSyntaxList {
					slow = slow.Parent
				}
				fast := getTouchingPropertyName_fast(file, pos)
				if fast.Kind == ast.KindJSDoc && slow.Kind == ast.KindIdentifier && slow.End() < pos {
					// The slow version (ported from Strada) has a bug where it can return an identifier
					// inside JSDoc where the position isn't actually touching the end of the identifier.
					return
				}
				slowToken, fastToken := toTokenInfo(slow), toTokenInfo(fast)
				assert.Equal(t, fastToken, slowToken)
			})
		}
	}
}

func TestGetTouchingPropertyName(t *testing.T) {
	t.Parallel()
	jstest.SkipIfNoNodeJS(t)
	repo.SkipIfNoTypeScriptSubmodule(t)
	for _, fileName := range testFiles {
		fileText, err := os.ReadFile(fileName)
		assert.NilError(t, err)
		file := parser.ParseSourceFile("file.ts", string(fileText), core.ScriptTargetLatest, scanner.JSDocParsingModeParseAll)
		positionCount := len(fileText)
		positions := make([]int, positionCount)
		for i := range positionCount {
			positions[i] = i * len(fileText) / positionCount
		}

		tsTokens := tsGetTouchingPropertyName(t, string(fileText), positions)
		t.Run(fileName, func(t *testing.T) {
			t.Parallel()
			for i, pos := range positions {
				t.Run(fmt.Sprintf("pos: %d", pos), func(t *testing.T) {
					t.Parallel()
					goToken := getTouchingPropertyName(file, pos)
					if goToken.Kind == ast.KindJSDocText && strings.HasPrefix(tsTokens[i].Kind, "JSDoc") {
						// Strada sometimes stored plain-text JSDoc comments as strings
						// on JSDoc nodes, whereas Corsa stores them as JSDocText nodes.
						// It's fine for Corsa to return a deeper, more specific node in
						// this case.
						return
					}
					assert.Equal(t, tsTokens[i], toTokenInfo(goToken))
				})
			}
		})
	}
}

type tokenInfo struct {
	Kind string `json:"kind"`
	Pos  int    `json:"pos"`
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
				};
			});
		};`

	info, err := jstest.EvalNodeScriptWithTS[[]tokenInfo](t, script, dir, core.Must(core.StringifyJson(positions, "", "")))
	assert.NilError(t, err)
	return info
}
