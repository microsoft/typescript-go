package ls

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"gotest.tools/v3/assert"
)

func TestGetOrganizeImportsActionsForTypeScriptKinds(t *testing.T) {
	t.Parallel()

	tests := []struct {
		requested    lsproto.CodeActionKind
		expected     lsproto.CodeActionKind
		expectedBase lsproto.CodeActionKind
	}{
		{lsproto.CodeActionKindSourceOrganizeImportsTs, lsproto.CodeActionKindSourceOrganizeImportsTs, lsproto.CodeActionKindSourceOrganizeImports},
		{lsproto.CodeActionKindSourceRemoveUnusedImportsTs, lsproto.CodeActionKindSourceRemoveUnusedImportsTs, lsproto.CodeActionKindSourceRemoveUnusedImports},
		{lsproto.CodeActionKindSourceSortImportsTs, lsproto.CodeActionKindSourceSortImportsTs, lsproto.CodeActionKindSourceSortImports},
	}

	for _, test := range tests {
		assert.DeepEqual(t, getOrganizeImportsActionsForKind(test.requested), []lsproto.CodeActionKind{test.expected})
		assert.Equal(t, getBaseOrganizeImportsKind(test.requested), test.expectedBase)
	}
}

func TestIsFixAllKindAcceptsTypeScriptKind(t *testing.T) {
	t.Parallel()

	assert.Assert(t, isFixAllKind(lsproto.CodeActionKindSourceFixAllTs))
}
