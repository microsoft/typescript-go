package ls_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil/projecttestutil"
	"gotest.tools/v3/assert"
)

func TestUntitledReferences(t *testing.T) {
	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	// First test the URI conversion functions to understand the issue
	untitledURI := lsproto.DocumentUri("untitled:Untitled-2")
	convertedFileName := ls.DocumentURIToFileName(untitledURI)
	t.Logf("URI '%s' converts to filename '%s'", untitledURI, convertedFileName)
	
	backToURI := ls.FileNameToDocumentURI(convertedFileName)
	t.Logf("Filename '%s' converts back to URI '%s'", convertedFileName, backToURI)
	
	if string(backToURI) != string(untitledURI) {
		t.Errorf("Round-trip conversion failed: '%s' -> '%s' -> '%s'", untitledURI, convertedFileName, backToURI)
	}

	// Create a simple test case with a regular file to simulate the issue
	testContent := `let x = 42;

x

x++;`

	regularFileName := "/Untitled-2.ts"
	
	// Set up the file system with a regular file
	files := map[string]any{
		regularFileName: testContent,
	}

	ctx := projecttestutil.WithRequestID(t.Context())
	service, done := createLanguageService(ctx, regularFileName, files)
	defer done()

	// Calculate position of 'x' on line 3 (zero-indexed line 2, character 0)
	position := 13 // After "let x = 42;\n\n"

	// Call ProvideReferences using the test method
	refs := service.TestProvideReferences(regularFileName, position)

	// Log the results
	t.Logf("Input file name: %s", regularFileName)
	t.Logf("Number of references found: %d", len(refs))
	for i, ref := range refs {
		t.Logf("Reference %d: URI=%s, Range=%+v", i+1, ref.Uri, ref.Range)
	}

	// We expect to find 3 references
	assert.Assert(t, len(refs) == 3, "Expected 3 references, got %d", len(refs))

	// Also test definition using ProvideDefinition
	uri := ls.FileNameToDocumentURI(regularFileName)
	lspPosition := lsproto.Position{Line: 2, Character: 0}
	definition, err := service.ProvideDefinition(t.Context(), uri, lspPosition)
	assert.NilError(t, err)
	if definition != nil && definition.Locations != nil {
		t.Logf("Definition found: %d locations", len(*definition.Locations))
		for i, loc := range *definition.Locations {
			t.Logf("Definition %d: URI=%s, Range=%+v", i+1, loc.Uri, loc.Range)
		}
	}
}