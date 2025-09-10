package sourcemap

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
)

// Test data based on TypeScript's declarationMapGoToDefinition test
const (
	// Source map content from the TypeScript test
	testSourceMapContent = `{"version":3,"file":"indexdef.d.ts","sourceRoot":"","sources":["index.ts"],"names":[],"mappings":"AAAA;IACI,MAAM,EAAE,MAAM,CAAC;IACf,UAAU,CAAC,QAAQ,EAAE,QAAQ,GAAG,IAAI;IACpC,WAAW;;;;;;;CAMd;AAED,MAAM,WAAW,QAAQ;IACrB,MAAM,EAAE,MAAM,CAAC;CAClB"}`

	// Declaration file content
	testDeclarationContent = `export declare class Foo {
    member: string;
    methodName(propName: SomeType): void;
    otherMethod(): {
        x: number;
        y?: undefined;
    } | {
        y: string;
        x?: undefined;
    };
}
export interface SomeType {
    member: number;
}
//# sourceMappingURL=indexdef.d.ts.map`

	// Original source file content
	testSourceContent = `export class Foo {
    member: string;
    methodName(propName: SomeType): void {}
    otherMethod() {
        if (Math.random() > 0.5) {
            return {x: 42};
        }
        return {y: "yes"};
    }
}

export interface SomeType {
    member: number;
}`
)

type testHost struct {
	files map[string]string
}

func (h *testHost) GetSource(fileName string) Source {
	content, exists := h.files[fileName]
	if !exists {
		return nil
	}
	return &testSourceFile{
		fileName: fileName,
		text:     content,
		lineMap:  core.ComputeLineStarts(content),
	}
}

func (h *testHost) GetCanonicalFileName(path string) string {
	return path // Case-sensitive for test
}

func (h *testHost) Log(text string) {
	// For testing, we can ignore logs or capture them
}

type testSourceFile struct {
	fileName string
	text     string
	lineMap  []core.TextPos
}

func (f *testSourceFile) FileName() string {
	return f.fileName
}

func (f *testSourceFile) Text() string {
	return f.text
}

func (f *testSourceFile) LineMap() []core.TextPos {
	return f.lineMap
}

func TestDocumentPositionMapper_GetSourcePosition(t *testing.T) {
	t.Parallel()
	host := &testHost{
		files: map[string]string{
			"/indexdef.d.ts": testDeclarationContent,
			"/index.ts":      testSourceContent,
		},
	}

	mapper := CreateDocumentPositionMapper(host, testSourceMapContent, "/indexdef.d.ts.map")

	// Test mapping from declaration file to source file
	// Position should be somewhere in the methodName declaration in the .d.ts file
	declLineStarts := core.ComputeLineStarts(testDeclarationContent)

	// Find the position of "methodName" in the declaration file
	// This is on line 2 (0-indexed), around character 4
	methodNamePosInDecl := int(declLineStarts[2]) + 4

	input := DocumentPosition{
		FileName: "/indexdef.d.ts",
		Pos:      core.TextPos(methodNamePosInDecl),
	}

	result := mapper.GetSourcePosition(input)

	// Should map to the source file
	if result.FileName != "/index.ts" {
		t.Errorf("Expected fileName to be '/index.ts', got '%s'", result.FileName)
	}

	// Should map to a position in the source file (we don't need exact position matching for this test)
	if result.Pos < 0 {
		t.Errorf("Expected positive position, got %d", result.Pos)
	}

	// Verify it's different from the input (actual mapping occurred)
	if result.FileName == input.FileName && result.Pos == input.Pos {
		t.Error("Expected mapping to change position, but got same position")
	}
}

func TestDocumentPositionMapper_NoMapping(t *testing.T) {
	t.Parallel()
	host := &testHost{
		files: map[string]string{},
	}

	// Test with empty source map
	mapper := CreateDocumentPositionMapper(host, `{"version":3,"file":"test.d.ts","sources":[],"mappings":""}`, "/test.d.ts.map")

	input := DocumentPosition{
		FileName: "/test.d.ts",
		Pos:      core.TextPos(10),
	}

	result := mapper.GetSourcePosition(input)

	// Should return unchanged input when no mappings exist
	if result.FileName != input.FileName || result.Pos != input.Pos {
		t.Errorf("Expected unchanged position, got fileName='%s' pos=%d", result.FileName, result.Pos)
	}
}

func TestDocumentPositionMapper_InvalidSourceMap(t *testing.T) {
	t.Parallel()
	host := &testHost{
		files: map[string]string{},
	}

	// Test with invalid JSON
	mapper := CreateDocumentPositionMapper(host, `invalid json`, "/test.d.ts.map")

	input := DocumentPosition{
		FileName: "/test.d.ts",
		Pos:      core.TextPos(10),
	}

	result := mapper.GetSourcePosition(input)

	// Should return unchanged input when source map is invalid
	if result.FileName != input.FileName || result.Pos != input.Pos {
		t.Errorf("Expected unchanged position with invalid source map, got fileName='%s' pos=%d", result.FileName, result.Pos)
	}
}

func TestCreateSourceMapper(t *testing.T) {
	t.Parallel()
	// Mock file reader
	fileReader := &testFileReader{
		files: map[string]string{
			"/indexdef.d.ts":     testDeclarationContent,
			"/indexdef.d.ts.map": testSourceMapContent,
			"/index.ts":          testSourceContent,
		},
	}

	// Mock host
	host := &testSourceMapperHost{}

	sourceMapper := CreateSourceMapper(host, fileReader)

	// Test mapping from declaration to source
	input := DocumentPosition{
		FileName: "/indexdef.d.ts",
		Pos:      core.TextPos(50), // Some position in the declaration file
	}

	result := sourceMapper.TryGetSourcePosition(input)
	if result == nil {
		t.Error("Expected source position mapping, got nil")
	} else if result.FileName != "/index.ts" {
		t.Errorf("Expected mapping to '/index.ts', got '%s'", result.FileName)
	}
}

// Test helper types
type testFileReader struct {
	files map[string]string
}

func (r *testFileReader) ReadFile(path string) (string, bool) {
	content, exists := r.files[path]
	return content, exists
}

type testSourceMapperHost struct{}

func (h *testSourceMapperHost) GetSource(fileName string) Source {
	// This would normally read from files, but for testing we'll create mock data
	switch fileName {
	case "/indexdef.d.ts":
		return &testSourceFile{
			fileName: fileName,
			text:     testDeclarationContent,
			lineMap:  core.ComputeLineStarts(testDeclarationContent),
		}
	case "/index.ts":
		return &testSourceFile{
			fileName: fileName,
			text:     testSourceContent,
			lineMap:  core.ComputeLineStarts(testSourceContent),
		}
	}
	return nil
}

func (h *testSourceMapperHost) GetCanonicalFileName(path string) string {
	return path
}

func (h *testSourceMapperHost) Log(text string) {
	// Ignore for testing
}

func (h *testSourceMapperHost) UseCaseSensitiveFileNames() bool {
	return true
}

func (h *testSourceMapperHost) GetCurrentDirectory() string {
	return "/"
}

func (h *testSourceMapperHost) ReadFile(path string) (string, bool) {
	// Not used in this test
	return "", false
}

func (h *testSourceMapperHost) FileExists(path string) bool {
	// Not used in this test
	return false
}
