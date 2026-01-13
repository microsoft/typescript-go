package ls

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/parser"
)

// Test for issue: Crash on find-all-references on `this`
// This reproduces the panic when getPossibleSymbolReferencePositions is called
// with a container from one file while searching in another file.
// Without the fix, this test panics with "slice bounds out of range".
func TestGetPossibleSymbolReferencePositions_CrossFileContainer(t *testing.T) {
	t.Parallel()

	// Create two files with very different sizes
	// File 1: A large file with many statements
	largeFileText := `namespace LargeNamespace {
    export class LargeClass {
        // Many properties to create a large file
        private prop01: string = "value 01 with extra text to make it much longer for testing";
        private prop02: string = "value 02 with extra text to make it much longer for testing";
        private prop03: string = "value 03 with extra text to make it much longer for testing";
        private prop04: string = "value 04 with extra text to make it much longer for testing";
        private prop05: string = "value 05 with extra text to make it much longer for testing";
        private prop06: string = "value 06 with extra text to make it much longer for testing";
        private prop07: string = "value 07 with extra text to make it much longer for testing";
        private prop08: string = "value 08 with extra text to make it much longer for testing";
        private prop09: string = "value 09 with extra text to make it much longer for testing";
        private prop10: string = "value 10 with extra text to make it much longer for testing";
        private prop11: string = "value 11 with extra text to make it much longer for testing";
        private prop12: string = "value 12 with extra text to make it much longer for testing";
        
        constructor() {
            this.prop01 = "initialized";
            this.prop02 = "initialized";
            this.prop03 = "initialized";
        }
        
        method() {
            console.log("Some method with this reference");
            console.log("Adding more lines");
            console.log("To push the position forward");
            return this;
        }
    }
}`

	largeFile := parser.ParseSourceFile(ast.SourceFileParseOptions{
		FileName: "/largeFile.ts",
		Path:     "/largeFile.ts",
	}, largeFileText, core.ScriptKindTS)

	// File 2: A very small file
	smallFileText := `function f() { return this; }`

	smallFile := parser.ParseSourceFile(ast.SourceFileParseOptions{
		FileName: "/smallFile.ts",
		Path:     "/smallFile.ts",
	}, smallFileText, core.ScriptKindTS)

	// Get any child node from the large file to use as a container
	// The key is that this container belongs to largeFile, not smallFile
	var containerFromLargeFile *ast.Node
	largeFile.AsNode().ForEachChild(func(child *ast.Node) bool {
		containerFromLargeFile = child
		return false // Stop after first child
	})

	if containerFromLargeFile == nil {
		t.Fatal("Could not find any node in large file")
	}

	// The crucial test: the container is from largeFile, but we're searching smallFile
	// Without the fix, this could cause issues if container.Pos() or container.End() 
	// have values that are invalid for smallFile's text length
	smallFileLength := len(smallFileText)
	t.Logf("Using container from large file to search small file")
	t.Logf("Small file length: %d", smallFileLength)
	
	// The bug occurs when container.End() > len(smallFileText) OR
	// when container.Pos() > len(smallFileText)
	// Let's test both scenarios by calling the function
	
	// Without the fix, this call could panic depending on the container's position/end values
	// The fix ensures that if the container is from a different file, we use the entire source file
	positions := getPossibleSymbolReferencePositions(smallFile, "this", containerFromLargeFile)

	// The function should handle this gracefully
	// We expect to find "this" in the small file
	if len(positions) == 0 {
		t.Error("Expected to find at least one reference to 'this' in small file")
	}

	// Verify positions are valid for the small file
	for _, pos := range positions {
		if pos < 0 || pos > smallFileLength {
			t.Errorf("Position %d is out of bounds for small file (length %d)", pos, smallFileLength)
		}
	}
	
	t.Logf("Successfully found %d positions without panicking", len(positions))
}