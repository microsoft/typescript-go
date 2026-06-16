package ast_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/testutil/parsetestutil"
	"gotest.tools/v3/assert"
)

func TestGetLineStarts(t *testing.T) {
	t.Parallel()
	text := `a
	b
	c`
	sourceFile := parsetestutil.ParseTypeScript(text, false)
	expectedStarts := []core.TextPos{0, 2, 5}
	actualStarts := sourceFile.GetLineStarts()
	assert.DeepEqual(t, expectedStarts, actualStarts)
}

func TestGetLineAndCharacterOfPosition(t *testing.T) {
	t.Parallel()

	data := []struct {
		title string
		input string
		jsx   bool
		position core.TextPos
		expectedLine int
		expectedCharacter int
	}{
		{
			title: "Empty file search",
			input: ``,
			position: 10,
			expectedLine: 0,
			expectedCharacter: 10,
		},
		{
			title: "Single line search",
			input: `function foo() {`,
			position: 10,
			expectedLine: 0,
			expectedCharacter: 10,
		},
		{
			title: "Multiple line search",
			input: `a
b`,
			position: 2,
			expectedLine: 1,
			expectedCharacter: 0,
		},
	}

	for _, rec := range data {
		t.Run("SourceFile "+rec.title, func(t *testing.T) {
			t.Parallel()

			sourceFile := parsetestutil.ParseTypeScript(rec.input, false)
			lineAndCharacter := sourceFile.GetLineAndCharacterOfPosition(rec.position)
			assert.Equal(t, lineAndCharacter.Line, rec.expectedLine)
			assert.Equal(t, lineAndCharacter.Character, rec.expectedCharacter)
		})
	}
}