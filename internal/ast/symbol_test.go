package ast_test

import (
	"testing"
	"unicode/utf8"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/json"
	"gotest.tools/v3/assert"
)

func TestSymbolNameEncoding(t *testing.T) {
	t.Parallel()

	internalName := ast.InternalSymbolNameCall
	userName := ast.EscapeLeadingUnderscores("__call")

	assert.Assert(t, utf8.ValidString(string(internalName)))
	assert.Assert(t, internalName != userName)
	assert.Equal(t, internalName.EscapedText(), "__call")
	assert.Equal(t, userName.EscapedText(), "___call")
	assert.Equal(t, ast.UnescapeLeadingUnderscores(userName), "__call")

	encoded, err := json.Marshal([]ast.SymbolNameKey{internalName, userName})
	assert.NilError(t, err)
	var decoded []string
	assert.NilError(t, json.Unmarshal(encoded, &decoded))
	assert.DeepEqual(t, decoded, []string{string(internalName), string(userName)})
}
