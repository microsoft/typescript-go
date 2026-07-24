package autoimport

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"gotest.tools/v3/assert"
)

func TestIsUnusableNameDistinguishesInternalNames(t *testing.T) {
	t.Parallel()

	assert.Assert(t, isUnusableName(ast.InternalSymbolNameExportStar))
	assert.Assert(t, !isUnusableName(ast.EscapeLeadingUnderscores("__export")))
}
