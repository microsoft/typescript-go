package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestCompletionNoCrashImportStatementLocationParentNil(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	emptyCommitChars := []string{}

	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "import equals newline then statement",
			content: "import x =/*1*/\nclass Foo {}",
		},
		{
			name:    "import equals trailing space newline then statement",
			content: "import x = /*1*/\nclass Foo {}",
		},
		{
			name:    "import equals newline only",
			content: "import x =/*1*/\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer testutil.RecoverAndFail(t, "Panic on fourslash test")
			f, done := fourslash.NewFourslash(t, nil /*capabilities*/, tt.content)
			defer done()
			f.VerifyCompletions(t, "1", &fourslash.CompletionsExpectedList{
				IsIncomplete: false,
				ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
					CommitCharacters: &emptyCommitChars,
				},
				Items: &fourslash.CompletionsExpectedItems{
					Includes: []fourslash.CompletionsExpectedItem{
						"type",
					},
				},
			})
		})
	}
}
